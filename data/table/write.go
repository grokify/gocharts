package table

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/pkg/errors"
)

// WriteCSV writes the table as a CSV.
func WriteCSV(path string, t *TableData) error {
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
	tbl := NewTableData()
	tbl.Columns = cols
	tbl.Records = records
	return tbl.WriteCSV(filename)
}

// WriteXLSX writes a table as an Excel XLSX file.
func WriteXLSX(path string, tbls ...*TableData) error {
	tfs := []*TableFormatter{}
	for _, tbl := range tbls {
		tfs = append(tfs,
			&TableFormatter{
				Table:     tbl,
				Formatter: FormatStrings,
			},
		)
	}
	return WriteXLSXFormatted(path, tfs...)
}

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

type TableFormatter struct {
	Table     *TableData
	Formatter func(val string, colIdx uint) (interface{}, error)
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
		for y, row := range t.Records {
			for x, cellValue := range row {
				cellLocation := CoordinatesToSheetLocation(uint32(x), uint32(y+rowBase))
				formattedVal, err := tf.Formatter(cellValue, uint(x))
				if err != nil {
					return errors.Wrap(err, "WriteXLSXFormatted.Error.FormatCellValue")
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

type jsonRecords struct {
	Records []map[string]string `json:"records,omitempty"`
}

func (t *TableData) WriteJSON(path string, perm os.FileMode, jsonPrefix, jsonIndent string) error {
	out := jsonRecords{Records: t.ToSliceMSS()}
	fmt.Printf("TABLE.WRITEJSON [%v]\n", path)
	bytes, err := jsonutil.MarshalSimple(out, jsonPrefix, jsonIndent)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, perm)
}
