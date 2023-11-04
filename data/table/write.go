package table

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table/sheet"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/time/timeutil"
	excelize "github.com/xuri/excelize/v2"
)

const excelizeLinkTypeExternal = "External"

var (
	ErrSheetNameCollision  = errors.New("sheet name collision")
	ErrTablesCannotBeEmpty = errors.New("tables cannot be empty")
	rxURLHTTPOrHTTPS       = regexp.MustCompile(`^(?i)https?://.`)
)

// WriteCSVSimple writes a file with cols and rows data.
func WriteCSVSimple(cols []string, rows [][]string, filename string) error {
	tbl := NewTable("")
	tbl.Columns = cols
	tbl.Rows = rows
	return tbl.WriteCSV(filename)
}

// FormatterFunc returns a formatter function. A custom format func is returned if it is
// supplied and `FormatMap` is empty. If FormatMap is not empty, a function for it is
// returned.`
func (tbl *Table) FormatterFunc() func(val string, colIdx uint) (any, error) {
	if tbl.FormatMap == nil || len(tbl.FormatMap) == 0 {
		if tbl.FormatFunc != nil {
			return tbl.FormatFunc
		}
		// return format.FormatStrings
	}

	return func(val string, colIdx uint) (any, error) {
		fmtType, ok := tbl.FormatMap[int(colIdx)]
		if !ok || len(strings.TrimSpace(fmtType)) == 0 {
			if fmtType, ok = tbl.FormatMap[-1]; !ok {
				fmtType = ""
			}
		}
		switch strings.ToLower(strings.TrimSpace(fmtType)) {
		case FormatFloat:
			if strings.TrimSpace(val) == "" {
				return float64(0), nil
			} else if floatVal, err := strconv.ParseFloat(val, 64); err != nil {
				return val, err
			} else {
				return floatVal, nil
			}
		case FormatInt:
			if strings.TrimSpace(val) == "" {
				return int(0), nil
			}
			intVal, err := strconv.Atoi(val)
			if err != nil {
				floatVal, err2 := strconv.ParseFloat(val, 64)
				if err2 != nil {
					return val, err
				}
				return int(floatVal), nil
			}
			return intVal, nil
		case FormatDate:
			if strings.TrimSpace(val) == "" {
				return "", nil // if date is not present, return an empty string.
			} else if dtVal, err := time.Parse(time.RFC3339, val); err != nil {
				return val, err
			} else {
				return dtVal.Format(timeutil.DateMDY), nil
			}
		case FormatTime:
			if strings.TrimSpace(val) == "" {
				return "", nil // if date is not present, return an empty string.
			} else if dtVal, err := time.Parse(time.RFC3339, val); err != nil {
				return val, err
			} else {
				return dtVal, nil
			}
		}
		return val, nil
	}
}

// WriteXLSX writes a table as an Excel XLSX file with row formatter option.
func WriteXLSX(path string, tbls []*Table) error {
	tables := []*Table{}
	for _, tbl := range tbls {
		if tbl != nil {
			tables = append(tables, tbl)
		}
	}
	if len(tables) == 0 {
		return ErrTablesCannotBeEmpty
	}

	sheetNames := map[string]int{} // track to avoid collisions and overwriting sheets

	f := excelize.NewFile()

	// Create a new sheet.
	sheetNum := 0
	for i, tbl := range tables {
		if tbl == nil {
			continue
		}
		sheetNum++
		sheetName := strings.TrimSpace(tbl.Name)
		if len(sheetName) == 0 {
			sheetName = fmt.Sprintf("Sheet%d", sheetNum)
		}
		if _, ok := sheetNames[sheetName]; ok {
			return errorsutil.Wrap(ErrSheetNameCollision, "sheet name collision for (%s)", sheetName)
		} else {
			sheetNames[sheetName]++
		}
		sheetIndex, err := f.NewSheet(sheetName)
		if err != nil {
			return errorsutil.Wrap(err, "excelize.File.NewSheet()")
		}
		// Set value of a cell.
		rowBase := 0
		if len(tbl.Columns) > 0 {
			rowBase++
			for i, cellValue := range tbl.Columns {
				cellLocation := sheet.CoordinatesToSheetLocation(uint(i), 0)
				err := f.SetCellValue(sheetName, cellLocation, cellValue)
				if err != nil {
					return err
				}
			}
		}
		fmtFunc := tbl.FormatterFunc()
		for y, row := range tbl.Rows {
			for x, cellValue := range row {
				cellLocation := sheet.CoordinatesToSheetLocation(uint(x), uint(y+rowBase))
				if fmtType, ok := tbl.FormatMap[x]; ok {
					if fmtType == FormatURL {
						txt, lnk := markdown.ParseLink(cellValue)
						txt = strings.TrimSpace(txt)
						lnk = strings.TrimSpace(lnk)
						if txt == "" && lnk != "" {
							txt = lnk
						}
						if txt != "" && lnk != "" {
							if err := f.SetCellValue(sheetName, cellLocation, txt); err != nil {
								return err
							}
							if err := f.SetCellHyperLink(sheetName, cellLocation, lnk, excelizeLinkTypeExternal); err != nil {
								return err
							}
							continue
						} else if rxURLHTTPOrHTTPS.MatchString(cellValue) {
							if err := f.SetCellValue(sheetName, cellLocation, cellValue); err != nil {
								return err
							}
							if err := f.SetCellHyperLink(sheetName, cellLocation, cellValue, excelizeLinkTypeExternal); err != nil {
								return err
							}
							continue
						}
					}
				}
				formattedVal, err := fmtFunc(cellValue, uint(x))
				if err != nil {
					return errorsutil.Wrap(err, "gocharts/data/tables/write.go/WriteXLSXFormatted.Error.FormatCellValue")
				}
				if err = f.SetCellValue(sheetName, cellLocation, formattedVal); err != nil {
					return err
				}
				if tbl.FormatAutoLink {
					if rxURLHTTPOrHTTPS.MatchString(cellValue) {
						err := f.SetCellHyperLink(sheetName, cellLocation, cellValue, excelizeLinkTypeExternal)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		// Set active sheet of the workbook.
		if i == 0 {
			f.SetActiveSheet(sheetIndex)
		}
	}

	// Delete default sheet.
	err := f.DeleteSheet(f.GetSheetName(0))
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.DeleteSheet()")
	}

	// Save xlsx file by the given path.
	return f.SaveAs(path)
}

type SheetData struct {
	SheetName string
	Rows      [][]any
}

func WriteXLSXInterface(filename string, sheetdatas ...SheetData) error {
	f := excelize.NewFile()
	// Delete default sheet.
	shtIndex, err := f.GetSheetIndex("Sheet1")
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.GetSheetIndex()")
	}
	err = f.DeleteSheet(f.GetSheetName(shtIndex))
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.DeleteSheet()")
	}
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return errorsutil.Wrap(err, "excelize.File.DeleteSheet()")
	}
	// Create a new sheet.
	for i, sheetdata := range sheetdatas {
		sheetname := strings.TrimSpace(sheetdata.SheetName)
		if len(sheetname) == 0 {
			sheetname = fmt.Sprintf("Sheet%d", i+1)
		}
		index, err := f.NewSheet(sheetname)
		if err != nil {
			return errorsutil.Wrap(err, "excelize.File.NewSheet()")
		}
		for y, row := range sheetdata.Rows {
			for x, cellValue := range row {
				cellLocation := sheet.CoordinatesToSheetLocation(uint(x), uint(y))
				err := f.SetCellValue(sheetname, cellLocation, cellValue)
				if err != nil {
					return err
				}
			}
		}
		// Set active sheet of the workbook.
		if i == 0 {
			f.SetActiveSheet(index)
		}
	}
	// Save xlsx file by the given path.
	return f.SaveAs(filename)
}

type jsonRecords struct {
	Records []map[string]string `json:"records,omitempty"`
}

func (tbl *Table) WriteJSON(path string, perm os.FileMode, jsonPrefix, jsonIndent string) error {
	out := jsonRecords{Records: tbl.ToSliceMSS()}
	fmt.Printf("TABLE.WRITEJSON [%v]\n", path)
	bytes, err := jsonutil.MarshalSimple(out, jsonPrefix, jsonIndent)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, perm)
}

/*
func WriteXLSXMapStringInt(filename, sheetname, colNameKey, colNameVal string, m map[string]int) error {
	tbl := NewTable("")
	if strings.TrimSpace(colNameKey) == "" {
		colNameKey = "Key"
	}
	if strings.TrimSpace(colNameVal) == "" {
		colNameVal = "Count"
	}
	tbl.Columns = []string{colNameKey, colNameVal}
	tbl.FormatMap = map[int]string{0: FormatString, 1: FormatInt}
	for k, v := range m {
		tbl.Rows = append(tbl.Rows, []string{k, strconv.Itoa(v)})
	}
	return tbl.WriteXLSX(filename, sheetname)
}
*/

func NewTableMapStringInt(tableName, colNameKey, colNameVal string, m map[string]int) Table {
	tbl := NewTable(tableName)
	if strings.TrimSpace(colNameKey) == "" {
		colNameKey = "Key"
	}
	if strings.TrimSpace(colNameVal) == "" {
		colNameVal = "Value"
	}
	tbl.Columns = []string{colNameKey, colNameVal}
	tbl.FormatMap = map[int]string{0: FormatString, 1: FormatInt}
	for k, v := range m {
		tbl.Rows = append(tbl.Rows, []string{k, strconv.Itoa(v)})
	}
	return tbl
}
