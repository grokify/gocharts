package table

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/grokify/gotilla/encoding/csvutil"
	"github.com/grokify/gotilla/encoding/jsonutil"
	"github.com/grokify/gotilla/type/stringsutil"
)

func NewTableDataFilesSimple(filenames []string, sep string, hasHeader, trimSpace bool) (TableData, error) {
	tbl := NewTableData()
	for i, filename := range filenames {
		filename = strings.TrimSpace(filename)
		if len(filename) == 0 {
			continue
		}
		tblx, err := NewTableDataFileSimple(filename, sep, hasHeader, trimSpace)
		if err != nil {
			return tbl, err
		}
		if len(tbl.Columns) == 0 {
			tbl.Columns = tblx.Columns
		} else {
			curCols := strings.Join(tbl.Columns, ",")
			nowCols := strings.Join(tblx.Columns, ",")
			if curCols != nowCols {
				if i == 0 {
					// if len(tbl.Columns) > 0, i has to be > 0
					panic("E_BAD_FILE_COUNTER_TABLE_COLUMNS")
				}
				return tbl, fmt.Errorf("CSV table definition mismatch [%s] AND [%s] for FILES [%s]",
					curCols, nowCols,
					filenames[i-1]+","+filename)
			}
		}
		if len(tblx.Records) > 0 {
			tbl.Records = append(tbl.Records, tblx.Records...)
		}
	}
	return tbl, nil
}

func NewTableDataFileSimple(path string, sep string, hasHeader, trimSpace bool) (TableData, error) {
	tbl := NewTableData()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return tbl, err
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if trimSpace {
			line = strings.TrimSpace(line)
		}
		parts := strings.Split(line, sep)
		parts = stringsutil.SliceTrimSpace(parts, false)
		if hasHeader && i == 0 {
			tbl.Columns = parts
		} else {
			tbl.Records = append(tbl.Records, parts)
		}
	}
	return tbl, nil
}

// NewTableDataFileCSV reads in a CSV file and returns a TableData struct.
func NewTableDataFileCSV(path string, comma rune, stripBom bool) (TableData, error) {
	tbl := NewTableData()
	csvReader, f, err := csvutil.NewReader(path, comma, stripBom)
	if err != nil {
		return tbl, err
	}
	defer f.Close()
	if DebugReadCSV {
		i := -1
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return tbl, err
			}
			i++
			if i == 0 {
				tbl.Columns = line
				continue
			}
			tbl.Records = append(tbl.Records, line)
			if i > 2500 {
				fmt.Printf("[%v] %v\n", i, strings.Join(line, ","))
			}
		}

	} else {
		lines, err := csvReader.ReadAll()
		if err != nil {
			return tbl, err
		}
		tbl.LoadMergedRows(lines)
	}
	return tbl, nil
}

func ReadMergeFilterCSVFiles(inPaths []string, outPath string, inComma rune, inStripBom bool, andFilter map[string]stringsutil.MatchInfo) (DocumentsSet, error) {
	//data := JsonRecordsInfo{Records: []map[string]string{}}
	data := NewDocumentsSet()

	for _, inPath := range inPaths {
		reader, inFile, err := csvutil.NewReader(inPath, inComma, inStripBom)
		if err != nil {
			return data, err
		}

		csvHeader := csvutil.CSVHeader{}
		j := -1

		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return data, err
			}
			j += 1

			if j == 0 {
				csvHeader.Columns = line
				continue
			}
			match, err := csvHeader.RecordMatch(line, andFilter)
			if err != nil {
				return data, err
			}
			if !match {
				continue
			}

			mss := csvHeader.RecordToMSS(line)
			data.Documents = append(data.Documents, mss)
		}
		err = inFile.Close()
		if err != nil {
			return data, err
		}
	}
	data.Inflate()
	return data, nil
}

func MergeFilterCSVFilesToJSON(inPaths []string, outPath string, inComma rune, inStripBom bool, perm os.FileMode, andFilter map[string]stringsutil.MatchInfo) error {
	data, err := ReadMergeFilterCSVFiles(inPaths, outPath, inComma, inStripBom, andFilter)
	if err != nil {
		return err
	}
	bytes, err := jsonutil.MarshalSimple(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outPath, bytes, perm)
}

func ReadCSVFilesSingleColumnValuesString(files []string, sep string, hasHeader, trimSpace bool, col uint, condenseUniqueSort bool) ([]string, error) {
	values := []string{}
	for _, file := range files {
		fileValues, err := ReadCSVFileSingleColumnValuesString(
			file, sep, hasHeader, trimSpace, col, false)
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

func ReadCSVFileSingleColumnValuesString(filename, sep string, hasHeader, trimSpace bool, col uint, condenseUniqueSort bool) ([]string, error) {
	tbl, err := NewTableDataFileSimple(filename, sep, hasHeader, trimSpace)
	if err != nil {
		return []string{}, err
	}
	values := []string{}
	for _, row := range tbl.Records {
		if len(row) > int(col) {
			values = append(values, row[col])
		}
	}
	if condenseUniqueSort {
		values = stringsutil.SliceCondenseSpace(values, true, true)
	}
	return values, nil
}
