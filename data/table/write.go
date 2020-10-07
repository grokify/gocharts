package table

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/pkg/errors"
)

// WriteCSV writes the table as a CSV.
func WriteCSV(path string, t *Table) error {
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
	err = writer.WriteAll(t.Records)
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

// WriteCSVSimple writes a file with cols and records data.
func WriteCSVSimple(cols []string, records [][]string, filename string) error {
	tbl := NewTable()
	tbl.Columns = cols
	tbl.Records = records
	return tbl.WriteCSV(filename)
}

/*
// WriteXLSX writes a table as an Excel XLSX file.
func WriteXLSX(path string, tbls ...*TableData) error {
	tfs := []*TableFormatter{}
	for _, tbl := range tbls {
		tfs = append(tfs,
			&TableFormatter{
				Table:     tbl,
				Formatter: FormatStrings})
	}
	return WriteXLSXFormatted(path, tfs...)
}
*/

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

/*
type TableFormatter struct {
	Table     *TableData
	FormatMap map[int]string
	Formatter func(val string, colIdx uint) (interface{}, error)
}

func (tf *TableFormatter) FormatterFunc() func(val string, colIdx uint) (interface{}, error) {
	if tf.FormatMap == nil || len(tf.FormatMap) == 0 {
		if tf.Formatter == nil {
			return FormatStrings
		}
		return tf.Formatter
	}
	return func(val string, colIdx uint) (interface{}, error) {
		if fmtType, ok := tf.FormatMap[int(colIdx)]; ok {
			fmtType := strings.ToLower(strings.TrimSpace(fmtType))
			if len(fmtType) > 0 {
				switch fmtType {
				case "float":
					floatVal, err := strconv.ParseFloat(val, 64)
					if err != nil {
						return val, err
					}
					return floatVal, nil
				case "int":
					intVal, err := strconv.Atoi(val)
					if err != nil {
						return val, err
					}
					return intVal, nil
				case "time":
					dtVal, err := time.Parse(time.RFC3339, val)
					if err != nil {
						return val, err
					} else {
						return dtVal, nil
					}
				}
			}
		}
		return val, nil
	}
}

// WriteXLSXFormatted writes a table as an Excel XLSX file with
// row formatter option.
func WriteXLSXFormatted(path string, tbls ...*TableFormatter) error {
	f := excelize.NewFile()
	// Delete default sheet.
	f.DeleteSheet(f.GetSheetName(f.GetSheetIndex("Sheet1")))
	f.DeleteSheet("Sheet1")
	// Create a new sheet.
	sheetNum := 0
	for i, tf := range tbls {
		if tf == nil || tf.Table == nil {
			continue
		}
		t := tf.Table
		sheetNum++
		sheetname := strings.TrimSpace(t.Name)
		if len(sheetname) == 0 {
			sheetname = fmt.Sprintf("Sheet%d", sheetNum)
		}
		index := f.NewSheet(sheetname)
		// Set value of a cell.
		rowBase := 0
		if len(t.Columns) > 0 {
			rowBase++
			for i, cellValue := range t.Columns {
				cellLocation := CoordinatesToSheetLocation(uint32(i), 0)
				f.SetCellValue(sheetname, cellLocation, cellValue)
			}
		}
		fmtFunc := tf.FormatterFunc()
		for y, row := range t.Records {
			for x, cellValue := range row {
				cellLocation := CoordinatesToSheetLocation(uint32(x), uint32(y+rowBase))
				formattedVal, err := fmtFunc(cellValue, uint(x))
				if err != nil {
					return errors.Wrap(err, "gocharts/data/tables/write.go/WriteXLSXFormatted.Error.FormatCellValue")
				}
				f.SetCellValue(sheetname, cellLocation, formattedVal)
			}
		}
		// Set active sheet of the workbook.
		if i == 0 {
			f.SetActiveSheet(index)
		}
	}
	// Save xlsx file by the given path.
	return f.SaveAs(path)
}
*/

func (tbl *Table) FormatterFunc() func(val string, colIdx uint) (interface{}, error) {
	if tbl.FormatMap == nil || len(tbl.FormatMap) == 0 {
		if tbl.FormatFunc == nil {
			return FormatStrings
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
		fmtType = strings.ToLower(strings.TrimSpace(fmtType))
		if len(fmtType) > 0 {
			switch fmtType {
			case "float":
				floatVal, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return val, err
				}
				return floatVal, nil
			case "int":
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return val, err
				}
				return intVal, nil
			case "time":
				dtVal, err := time.Parse(time.RFC3339, val)
				if err != nil {
					return val, err
				}
				return dtVal, nil
			}
		}
		return val, nil
	}
}

// WriteXLSX writes a table as an Excel XLSX file with
// row formatter option.
func WriteXLSX(path string, tbls ...*Table) error {
	f := excelize.NewFile()
	// Delete default sheet.
	f.DeleteSheet(f.GetSheetName(f.GetSheetIndex("Sheet1")))
	f.DeleteSheet("Sheet1")
	// Create a new sheet.
	sheetNum := 0
	for i, tbl := range tbls {
		if tbl == nil {
			continue
		}
		//t := tf.Table
		sheetNum++
		sheetname := strings.TrimSpace(tbl.Name)
		if len(sheetname) == 0 {
			sheetname = fmt.Sprintf("Sheet%d", sheetNum)
		}
		index := f.NewSheet(sheetname)
		// Set value of a cell.
		rowBase := 0
		if len(tbl.Columns) > 0 {
			rowBase++
			for i, cellValue := range tbl.Columns {
				cellLocation := CoordinatesToSheetLocation(uint32(i), 0)
				f.SetCellValue(sheetname, cellLocation, cellValue)
			}
		}
		fmtFunc := tbl.FormatterFunc()
		for y, row := range tbl.Records {
			for x, cellValue := range row {
				cellLocation := CoordinatesToSheetLocation(uint32(x), uint32(y+rowBase))
				formattedVal, err := fmtFunc(cellValue, uint(x))
				if err != nil {
					return errors.Wrap(err, "gocharts/data/tables/write.go/WriteXLSXFormatted.Error.FormatCellValue")
				}
				f.SetCellValue(sheetname, cellLocation, formattedVal)
			}
		}
		// Set active sheet of the workbook.
		if i == 0 {
			f.SetActiveSheet(index)
		}
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
	f.DeleteSheet(f.GetSheetName(f.GetSheetIndex("Sheet1")))
	f.DeleteSheet("Sheet1")
	// Create a new sheet.
	for i, sheetdata := range sheetdatas {
		sheetname := strings.TrimSpace(sheetdata.SheetName)
		if len(sheetname) == 0 {
			sheetname = fmt.Sprintf("Sheet%d", i+1)
		}
		index := f.NewSheet(sheetname)
		for y, row := range sheetdata.Rows {
			for x, cellValue := range row {
				cellLocation := CoordinatesToSheetLocation(uint32(x), uint32(y))
				f.SetCellValue(sheetname, cellLocation, cellValue)
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

func (t *Table) WriteJSON(path string, perm os.FileMode, jsonPrefix, jsonIndent string) error {
	out := jsonRecords{Records: t.ToSliceMSS()}
	fmt.Printf("TABLE.WRITEJSON [%v]\n", path)
	bytes, err := jsonutil.MarshalSimple(out, jsonPrefix, jsonIndent)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, perm)
}
