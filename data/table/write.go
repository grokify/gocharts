package table

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gocharts/v2/data/table/format"
	"github.com/grokify/gocharts/v2/data/table/sheet"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	excelize "github.com/xuri/excelize/v2"
)

func writeCSV(path string, t *Table) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	if len(t.Columns) > 0 {
		err = writer.Write(t.Columns)
		if err != nil {
			return err
		}
	}
	err = writer.WriteAll(t.Rows)
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// WriteCSVSimple writes a file with cols and rows data.
func WriteCSVSimple(cols []string, rows [][]string, filename string) error {
	tbl := NewTable("")
	tbl.Columns = cols
	tbl.Rows = rows
	return tbl.WriteCSV(filename)
}

/*
func FormatStrings(val string, col uint) (interface{}, error) {
	return val, nil
}

func FormatStringAndInts(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		return val, nil
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatStringAndFloats(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		return val, nil
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatTimeAndInts(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt, nil
		}
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatMonthAndFloats(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt.Format(timeutil.ISO8601YM), nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatDateAndFloats(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt.Format(timeutil.DateMDY), nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}

func FormatTimeAndFloats(val string, colIdx uint) (interface{}, error) {
	if colIdx == 0 {
		dt, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return val, err
		} else {
			return dt, nil
		}
	}
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val, err
	}
	return num, nil
}
*/

func (tbl *Table) FormatterFunc() func(val string, colIdx uint) (interface{}, error) {
	if tbl.FormatMap == nil || len(tbl.FormatMap) == 0 {
		if tbl.FormatFunc == nil {
			return format.FormatStrings
		}
		return tbl.FormatFunc
	}
	return func(val string, colIdx uint) (interface{}, error) {
		fmtType, ok := tbl.FormatMap[int(colIdx)]
		if !ok || len(strings.TrimSpace(fmtType)) == 0 {
			fmtType, ok = tbl.FormatMap[-1]
			if !ok {
				fmtType = ""
			}
		}
		switch strings.ToLower(strings.TrimSpace(fmtType)) {
		case FormatFloat:
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return val, err
			}
			return floatVal, nil
		case FormatInt:
			intVal, err := strconv.Atoi(val)
			if err != nil {
				floatVal, err2 := strconv.ParseFloat(val, 64)
				if err2 != nil {
					return val, err
				}
				return int(floatVal), nil
			}
			return intVal, nil
		case FormatTime:
			dtVal, err := time.Parse(time.RFC3339, val)
			if err != nil {
				return val, err
			}
			return dtVal, nil
		}
		return val, nil
	}
}

var rxURLHTTPOrHTTPS = regexp.MustCompile(`^(?i)https?://.`)

// WriteXLSX writes a table as an Excel XLSX file with row formatter option.
func WriteXLSX(path string, tables ...*Table) error {
	f := excelize.NewFile()
	// Create a new sheet.
	sheetNum := 0
	for i, tbl := range tables {
		if tbl == nil {
			continue
		}
		sheetNum++
		sheetname := strings.TrimSpace(tbl.Name)
		if len(sheetname) == 0 {
			sheetname = fmt.Sprintf("Sheet%d", sheetNum)
		}
		index, err := f.NewSheet(sheetname)
		if err != nil {
			return errorsutil.Wrap(err, "excelize.File.NewShet()")
		}
		// Set value of a cell.
		rowBase := 0
		if len(tbl.Columns) > 0 {
			rowBase++
			for i, cellValue := range tbl.Columns {
				cellLocation := sheet.CoordinatesToSheetLocation(uint32(i), 0)
				err := f.SetCellValue(sheetname, cellLocation, cellValue)
				if err != nil {
					return err
				}
			}
		}
		fmtFunc := tbl.FormatterFunc()
		for y, row := range tbl.Rows {
			for x, cellValue := range row {
				cellLocation := sheet.CoordinatesToSheetLocation(uint32(x), uint32(y+rowBase))
				formattedVal, err := fmtFunc(cellValue, uint(x))
				if err != nil {
					return errorsutil.Wrap(err, "gocharts/data/tables/write.go/WriteXLSXFormatted.Error.FormatCellValue")
				}
				err = f.SetCellValue(sheetname, cellLocation, formattedVal)
				if err != nil {
					return err
				}
				if tbl.FormatAutoLink {
					if rxURLHTTPOrHTTPS.MatchString(cellValue) {
						err := f.SetCellHyperLink(sheetname, cellLocation, cellValue, "External")
						if err != nil {
							return err
						}
					}
				}
			}
		}
		// Set active sheet of the workbook.
		if i == 0 {
			f.SetActiveSheet(index)
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
	Rows      [][]interface{}
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
				cellLocation := sheet.CoordinatesToSheetLocation(uint32(x), uint32(y))
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
