package table

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/errors/errorsutil"
)

var debugReadCSV = false // should not need to use this.

// ReadFile reads one or more delimited files and returns a merged `Table` struct.
// An error is returned if the columns count differs between files.
func ReadFile(opts *ParseOptions, filenames ...string) (Table, error) {
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
			return tbl, fmt.Errorf("csv column mismatch earlier files cols1 [%s] cols2 [%s]",
				strings.Join(tbl.Columns, ","),
				strings.Join(tblx.Columns, ","))
		}
		tbl.Rows = append(tbl.Rows, tblx.Rows...)
	}
	return tbl, nil
}

func trimSpaceSliceSliceString(s [][]string) [][]string {
	for i, row := range s {
		for j, cell := range row {
			s[i][j] = strings.TrimSpace(cell)
		}
	}
	return s
}

func readSingleFile(opts *ParseOptions, filename string) (Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return NewTable(""), err
	}
	defer f.Close()
	tbl, err := ParseReadSeeker(opts, f)
	if err != nil {
		return tbl, err
	}
	return tbl, f.Close()
}

// ParseReadSeeker parses an `io.ReadSeeker` and returns a `Table` struct.
func ParseReadSeeker(opts *ParseOptions, rs io.ReadSeeker) (Table, error) {
	tbl := NewTable("")
	if opts == nil {
		opts = &ParseOptions{}
	}
	csvReader, err := csvutil.NewReader(rs, opts.CommaValue())
	if err != nil {
		return tbl, err
	}
	if debugReadCSV {
		i := -1
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return tbl, errorsutil.Wrapf(err, "error on `csvReader.Read()` row [%d]", i)
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
		csvReader.FieldsPerRecord = opts.FieldsPerRecord
		if opts.NoHeader {
			csvReader.FieldsPerRecord = -1
		}
		errorOutofBounds := true
		if csvReader.FieldsPerRecord < 0 {
			errorOutofBounds = false
		}
		lines, err := csvReader.ReadAll()
		if opts.TrimSpace {
			lines = trimSpaceSliceSliceString(lines)
		}
		if err != nil {
			return tbl, errorsutil.Wrap(err, "csv.Reader.ReadAll()")
		}
		if len(lines) == 0 {
			return tbl, errors.New("no content")
		}
		if !opts.NoHeader { // hasHeader
			if !opts.HasFilter() {
				tbl.LoadMergedRows(lines)
			} else {
				if len(opts.FilterColNames) > 0 {
					tbl.Columns = Columns(opts.FilterColNames)
				} else {
					cols := lines[0]
					wantColNames := []string{}
					for _, idx := range opts.FilterColIndices {
						if lenCols := len(cols); lenCols < 0 {
							panic("len cannot be negative")
						} else if idx < uint(lenCols) {
							wantColNames = append(wantColNames, cols[idx])
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
			if len(opts.FilterColIndices) == 0 {
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

// ParseOptions provides a set of configuraation parameters
// for parsing a CSV file. If an empty or nil `ParseOptions`
// is provided to `ReadFile`, default options for reading
// files will be used.
type ParseOptions struct {
	UseComma         bool
	Comma            rune
	NoHeader         bool // HasHeader is default (false)
	FieldsPerRecord  int
	FilterColNames   []string
	FilterColIndices []uint
	TrimSpace        bool
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

func (opts *ParseOptions) Filter(cols []string, rows [][]string, errorOutOfBounds bool) ([][]string, error) {
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
			} else if errorOutOfBounds {
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
*/

// Unmarshal is a convenience function to provide a simple interface to
// unmarshal table contents into any desired output.
func (tbl *Table) Unmarshal(funcRow func(row []string) error) error {
	for i, row := range tbl.Rows {
		err := funcRow(row)
		if err != nil {
			return errorsutil.Wrap(err, fmt.Sprintf("Error on Record Index [%d]", i))
		}
	}
	return nil
}

// ReadTableSimple is used for CSVs without standard escaping. For example,
// when `"` is used within a string, but is not used for encapsulating strings.
func ReadTableSimple(filename, sep string) (*Table, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t := NewTable("")

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), sep)
		if i == 0 {
			t.Columns = parts
		} else {
			t.Rows = append(t.Rows, parts)
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &t, nil
}
