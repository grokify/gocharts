package table

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/grokify/simplego/encoding/csvutil"
	"github.com/grokify/simplego/type/stringsutil"
	"github.com/pkg/errors"
)

var debugReadCSV = false // should not need to use this.

// ReadFiles reads in a list of delimited files and returns a merged `Table` struct.
// An error is returned if the columns count differs between files.
func ReadFiles(filenames []string, comma rune, hasHeader bool) (Table, error) {
	tbl := NewTable()
	for i, filename := range filenames {
		tblx, err := ReadFile(filename, comma, hasHeader)
		if err != nil {
			return tblx, err
		}
		if i > 0 && len(tbl.Columns) != len(tblx.Columns) {
			return tbl, fmt.Errorf("csv column count mismatch earlier files count [%d] file [%s] count [%d]",
				len(tbl.Columns), filename, len(tblx.Columns))
		}
	}
	return tbl, nil
}

// ReadFile reads in a delimited file and returns a `Table` struct.
func ReadFile(filename string, comma rune, hasHeader bool) (Table, error) {
	tbl := NewTable()
	csvReader, f, err := csvutil.NewReader(filename, comma, false)
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
			if i == 0 && hasHeader {
				tbl.Columns = line
				continue
			}
			tbl.Rows = append(tbl.Rows, line)
			if i > 2500 {
				fmt.Printf("[%v] %v\n", i, strings.Join(line, ","))
			}
		}
	} else {
		lines, err := csvReader.ReadAll()
		if err != nil {
			return tbl, err
		}
		byteOrderMarkAsString := string('\uFEFF')
		if len(lines) > 0 && len(lines[0]) > 0 &&
			strings.HasPrefix(lines[0][0], byteOrderMarkAsString) {
			lines[0][0] = strings.TrimPrefix(lines[0][0], byteOrderMarkAsString)
		}
		if hasHeader {
			tbl.LoadMergedRows(lines)
		} else {
			tbl.Rows = lines
		}
	}
	return tbl, nil
}

func ReadCSVFilesSingleColumnValuesString(files []string, sep rune, hasHeader bool, col uint, condenseUniqueSort bool) ([]string, error) {
	values := []string{}
	for _, file := range files {
		fileValues, err := ReadCSVFileSingleColumnValuesString(
			file, sep, hasHeader, col, false)
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

func ReadCSVFileSingleColumnValuesString(filename string, sep rune, hasHeader bool, col uint, condenseUniqueSort bool) ([]string, error) {
	tbl, err := ReadFile(filename, sep, hasHeader)
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

func ParseBytes(data []byte, delimiter rune, hasHeaderRow bool) (Table, error) {
	return ParseReader(bytes.NewReader(data), delimiter, hasHeaderRow)
}

func ParseReader(reader io.Reader, delimiter rune, hasHeaderRow bool) (Table, error) {
	tbl := NewTable()
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
