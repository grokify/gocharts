package documents

import (
	"io"
	"os"

	"github.com/grokify/mogo/encoding/csvutil"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/type/stringsutil"
)

func ReadMergeFilterCSVFiles(inPaths []string, outPath string, inComma rune, andFilter map[string]stringsutil.MatchInfo) (DocumentsSet, error) {
	//data := JsonRecordsInfo{Records: []map[string]string{}}
	data := NewDocumentsSet()

	for _, inPath := range inPaths {
		reader, inFile, err := csvutil.NewReaderFile(inPath, inComma)
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
			j++

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

func MergeFilterCSVFilesToJSON(inPaths []string, outPath string, inComma rune, perm os.FileMode, andFilter map[string]stringsutil.MatchInfo) error {
	data, err := ReadMergeFilterCSVFiles(inPaths, outPath, inComma, andFilter)
	if err != nil {
		return err
	}
	bytes, err := jsonutil.MarshalSimple(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, bytes, perm)
}
