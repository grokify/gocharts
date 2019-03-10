package table

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/encoding/jsonutil"
)

// TableData is useful for working on CSV data
type TableData struct {
	Columns []string
	Records [][]string
}

func NewTableData() TableData {
	return TableData{
		Columns: []string{},
		Records: [][]string{}}
}

// NewTableDataFileCSV reads in a CSV file and returns a TableData struct.
func NewTableDataFileCSV(path string, comma rune, stripBom bool) (TableData, error) {
	tbl := TableData{}
	csvReader, f, err := csvutil.NewReader(path, comma, stripBom)
	if err != nil {
		return tbl, err
	}
	defer f.Close()
	lines, err := csvReader.ReadAll()
	if err != nil {
		return tbl, err
	}
	tbl.LoadMergedRows(lines)
	return tbl, nil
}

// LoadMergedRows is used to load data from `[][]string` sources
// like csv.ReadAll()
func (t *TableData) LoadMergedRows(data [][]string) {
	if len(data) == 0 {
		return
	}
	t.Columns = data[0]
	if len(data) > 1 {
		t.Records = data[1:]
	}
}

func (t *TableData) ColumnIndex(wantCol string) int {
	for i, col := range t.Columns {
		if col == wantCol {
			return i
		}
	}
	return -1
}

func (t *TableData) ColumnsValuesDistinct(wantCols []string, stripSpace bool) (map[string]int, error) {
	data := map[string]int{}
	wantIdxs := []int{}
	maxIdx := -1
	for _, wantCol := range wantCols {
		wantIdx := t.ColumnIndex(wantCol)
		if wantIdx < 0 {
			return data, fmt.Errorf("Column Not Found [%v]", wantCol)
		}
		wantIdxs = append(wantIdxs, wantIdx)
		if wantIdx > maxIdx {
			maxIdx = wantIdx
		}
	}
	for _, rec := range t.Records {
		if len(rec) > maxIdx {
			vals := []string{}
			for _, wantIdx := range wantIdxs {
				val := rec[wantIdx]
				if stripSpace {
					val = strings.TrimSpace(val)
				}
				vals = append(vals, val)
			}
			valsStr := strings.Join(vals, " ")
			_, ok := data[valsStr]
			if !ok {
				data[valsStr] = 0
			}
			data[valsStr] += 1
		}
	}
	return data, nil
}

func (t *TableData) ColumnValuesDistinct(wantCol string) (map[string]int, error) {
	data := map[string]int{}
	idx := t.ColumnIndex(wantCol)
	if idx < 0 {
		return data, fmt.Errorf("Column Not Found [%v]", wantCol)
	}

	for _, rec := range t.Records {
		if len(rec) > idx {
			val := rec[idx]
			_, ok := data[val]
			if !ok {
				data[val] = 0
			}
			data[val] += 1
		}
	}
	return data, nil
}

func (t *TableData) ColumnValuesMinMax(wantCol string) (string, string, error) {
	vals, err := t.ColumnValuesDistinct(wantCol)
	if err != nil {
		return "", "", err
	}
	if len(vals) == 0 {
		return "", "", errors.New("No Values Found")
	}

	arr := []string{}
	for val := range vals {
		arr = append(arr, val)
	}

	sort.Strings(arr)
	return arr[0], arr[len(arr)-1], nil
}

func (t *TableData) RecordValue(wantCol string, record []string) (string, error) {
	idx := t.ColumnIndex(wantCol)
	if idx < 0 {
		return "", fmt.Errorf("Column Not Found [%v]", wantCol)
	}
	if idx >= len(record) {
		return "", fmt.Errorf("Record does not have enough columns [%v]", idx+1)
	}
	return record[idx], nil
}

func (t *TableData) RecordValueOrEmpty(wantCol string, record []string) string {
	val, err := t.RecordValue(wantCol, record)
	if err != nil {
		return ""
	}
	return val
}

func (t *TableData) NewTableFiltered(wantColValues map[string]string) (TableData, error) {
	t2 := TableData{Columns: t.Columns}
	records, err := t.FilterRecords(wantColValues)
	if err != nil {
		return t2, err
	}
	t2.Records = records
	return t2, nil
}

func (t *TableData) FilterRecords(wantColValues map[string]string) ([][]string, error) {
	data := [][]string{}
	wantColIndexes := map[string]int{}
	maxIdx := -1
	for wantColName := range wantColValues {
		wantColIdx := t.ColumnIndex(wantColName)
		if wantColIdx < 0 {
			return data, fmt.Errorf("Column Not Found [%v]", wantColName)
		}
		if wantColIdx > maxIdx {
			maxIdx = wantColIdx
		}
		wantColIndexes[wantColName] = wantColIdx
	}
RECORDS:
	for _, rec := range t.Records {
		if len(rec) > maxIdx {
			for wantColName, wantColIdx := range wantColIndexes {
				colValue := rec[wantColIdx]
				wantColValue, ok := wantColValues[wantColName]
				if !ok {
					return data, fmt.Errorf("Column Name [%v] has no desired value", wantColName)
				}
				if colValue != wantColValue {
					continue RECORDS
				}
			}
			data = append(data, rec)
		}
	}
	return data, nil
}

func (t *TableData) WriteCSV(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(t.Columns)
	if err != nil {
		return err
	}
	err = writer.WriteAll(t.Records)
	if err != nil {
		return err
	}
	writer.Flush()
	return writer.Error()
}

func (t *TableData) RecordToMSS(record []string) map[string]string {
	mss := map[string]string{}
	for i, key := range t.Columns {
		if i < len(t.Columns) {
			mss[key] = record[i]
		}
	}
	return mss
}

func (t *TableData) ToSliceMSS() []map[string]string {
	slice := []map[string]string{}
	for _, rec := range t.Records {
		slice = append(slice, t.RecordToMSS(rec))
	}
	return slice
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
