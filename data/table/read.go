package table

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/grokify/simplego/encoding/csvutil"
	"github.com/pkg/errors"
)

var debugReadCSV = false // should not need to use this.

// ReadFiles reads in a list of delimited files and returns a merged `Table` struct.
// An error is returned if the columns count differs between files.
func ReadFiles(opts *ParseOptions, filenames ...string) (Table, error) {
	tbl := NewTable("")
	if len(filenames) == 0 {
		return tbl, errors.New("no filenames provided")
	}
	for i, filename := range filenames {
		tblx, err := readSingleFile(opts, filename)
		if err != nil || len(filenames) == 1 {
			return tblx, err
		}
		if i == 0 {
			tbl = tblx
			continue
		} else if !tbl.Columns.Equal(tblx.Columns) {
			return tbl, fmt.Errorf("csv column count mismatch earlier files count [%d] file [%s] count [%d]",
				len(tbl.Columns), filename, len(tblx.Columns))
		}
		tbl.Rows = append(tbl.Rows, tblx.Rows...)
	}
	return tbl, nil
}

// ReadFile reads in a delimited file and returns a `Table` struct.
func readSingleFile(opts *ParseOptions, filename string) (Table, error) {
	tbl := NewTable("")
	comma := ','
	if opts != nil {
		comma = opts.CommaValue()
	}
	csvReader, f, err := csvutil.NewReader(filename, comma)
	if err != nil {
		return tbl, err
	}
	defer f.Close()
	if debugReadCSV {
		i := -1
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return tbl, err
			}
			i++
			if i == 0 && (opts == nil || !opts.NoHeader) {
				tbl.Columns = line
				continue
			}
			tbl.Rows = append(tbl.Rows, line)
			if i > 2500 {
				fmt.Printf("[%v] %v\n", i, strings.Join(line, ","))
			}
		}
	} else {
		errorOutofBounds := true
		if opts != nil {
			csvReader.FieldsPerRecord = opts.FieldsPerRecord
			if csvReader.FieldsPerRecord < 0 {
				errorOutofBounds = false
			}
		}
		lines, err := csvReader.ReadAll()
		if err != nil {
			return tbl, errors.Wrap(err, "csv.Reader.ReadAll()")
		}
		if len(lines) == 0 {
			return tbl, errors.New("no content")
		}
		if len(lines) > 0 && len(lines[0]) > 0 {
			lines[0][0] = trimUTF8ByteOrderMarkString(lines[0][0])
		}
		if opts == nil || !opts.NoHeader {
			if opts == nil || !opts.HasFilter() {
				tbl.LoadMergedRows(lines)
			} else {
				if len(opts.FilterColNames) > 0 {
					tbl.Columns = Columns(opts.FilterColNames)
				} else {
					cols := lines[0]
					wantColNames := []string{}
					for _, idx := range opts.FilterColIndices {
						if int(idx) < len(cols) {
							wantColNames = append(wantColNames, cols[int(idx)])
						} else {
							return tbl, fmt.Errorf("want column index not found [%d]", idx)
						}
					}
					tbl.Columns = wantColNames
				}

				rows, err := opts.Filter(lines[0], lines[1:], errorOutofBounds)
				if err != nil {
					return tbl, err
				}
				tbl.Rows = rows
			}
		} else {
			if opts == nil || len(opts.FilterColIndices) == 0 {
				tbl.Rows = lines
			} else {
				rows, err := opts.Filter([]string{}, lines, errorOutofBounds)
				if err != nil {
					return tbl, err
				}
				tbl.Rows = rows
			}
		}
	}
	return tbl, nil
}

func trimUTF8ByteOrderMarkString(s string) string {
	byteOrderMarkAsString := string('\uFEFF')
	if strings.HasPrefix(s, byteOrderMarkAsString) {
		return strings.TrimPrefix(s, byteOrderMarkAsString)
	}
	return s
}

type ParseOptions struct {
	UseComma         bool
	Comma            rune
	NoHeader         bool // HasHeader is default (false)
	FieldsPerRecord  int
	FilterColNames   []string
	FilterColIndices []uint
}

func (opts *ParseOptions) CommaValue() rune {
	if opts.UseComma {
		return opts.Comma
	}
	return ','
}

func (opts *ParseOptions) HasFilter() bool {
	if len(opts.FilterColNames) > 0 || len(opts.FilterColIndices) > 0 {
		return true
	}
	return false
}

func (opts *ParseOptions) Filter(cols []string, rows [][]string, errorOutofBounds bool) ([][]string, error) {
	newRows := [][]string{}
	indices := opts.FilterColIndices

	if !opts.HasFilter() {
		return rows, nil
	}

	if len(opts.FilterColIndices) == 0 &&
		len(opts.FilterColNames) > 0 && len(cols) > 0 {
		colsPlus := Columns(cols)
		indicesTry := []uint{}
		wantColNamesNotFound := []string{}
		for _, wantColName := range opts.FilterColNames {
			idx := colsPlus.Index(wantColName)
			if idx < 0 {
				wantColNamesNotFound = append(wantColNamesNotFound, wantColName)
			} else {
				indicesTry = append(indicesTry, uint(idx))
			}
		}
		if len(wantColNamesNotFound) > 0 {
			return newRows, fmt.Errorf(
				"filter columns not found [%s]",
				strings.Join(wantColNamesNotFound, ","))
		}
		indices = indicesTry
	}

	if len(indices) == 0 {
		return newRows, fmt.Errorf(
			"no colIndices or row match: colNameFilter [%s] colNames [%s]",
			strings.Join(opts.FilterColNames, ","),
			strings.Join(cols, ","))
	}

	for _, row := range rows {
		newRow := []string{}
		for _, idx := range indices {
			if int(idx) <= len(row) {
				newRow = append(newRow, row[int(idx)])
			} else if errorOutofBounds {
				return newRows, fmt.Errorf("desired index out of bounds: index[%d] row len [%d]", idx, len(row))
			} else {
				newRow = append(newRow, "")
			}
		}
		newRows = append(newRows, newRow)
	}

	return newRows, nil
}

/*
func ReadCSVFilesSingleColumnValuesString(files []string, col uint, condenseUniqueSort bool) ([]string, error) {
	values := []string{}
	for _, file := range files {
		fileValues, err := ReadCSVFileSingleColumnValuesString(
			file, col, false)
		if err != nil {
			return values, err
		}
		values = append(values, fileValues...)
	}
	if condenseUniqueSort {
		values = stringsutil.SliceCondenseSpace(values, true, true)
	}
	return values, nil
}

func ReadCSVFileSingleColumnValuesString(filename string, col uint, condenseUniqueSort bool) ([]string, error) {
	tbl, err := ReadFile(filename, nil)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for _, row := range tbl.Rows {
		if len(row) > int(col) {
			values = append(values, row[col])
		}
	}
	if condenseUniqueSort {
		values = stringsutil.SliceCondenseSpace(values, true, true)
	}
	return values, nil
}
*/

func ParseBytes(data []byte, delimiter rune, hasHeaderRow bool) (Table, error) {
	return ParseReader(bytes.NewReader(data), delimiter, hasHeaderRow)
}

func ParseReader(reader io.Reader, delimiter rune, hasHeaderRow bool) (Table, error) {
	tbl := NewTable("")
	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimiter
	csvReader.TrimLeadingSpace = true
	idx := -1
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return tbl, err
		}
		idx++
		if idx == 0 && hasHeaderRow {
			tbl.Columns = row
			continue
		}
		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl, nil
}

// Unmarshal is a convenience function to provide a simple interface to
// unmarshal table contents into any desired output.
func (tbl *Table) Unmarshal(funcRow func(row []string) error) error {
	for i, row := range tbl.Rows {
		err := funcRow(row)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Error on Record Index [%d]", i))
		}
	}
	return nil
}
