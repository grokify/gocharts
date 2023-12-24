package histogram

import (
	"errors"
	"strconv"
	"strings"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
	"golang.org/x/exp/slices"
)

// AddMap provides a helper function to automatically create url encoded string keys.
// This can be used with `TableMap` to generate tables with arbitrary columns easily.
func (hist *Histogram) AddMap(binMap map[string]string, binCount int) {
	m := maputil.MapStringString(binMap)
	key := m.Encode()
	hist.Add(key, binCount)
}

// MapKeySplit returns a new `HistogramSet` where the supplied key is the HistogramSet key.
func (hist *Histogram) MapKeySplit(mapKey string, mapValIncl []string) (*HistogramSet, error) {
	hs := NewHistogramSet(mapKey)
	mapValInclMap := map[string]int{}
	for _, k := range mapValIncl {
		mapValInclMap[k]++
	}
	/*
		if 1 == 0 {
			mk, err := hist.MapKeys()
			logutil.FatalErr(err)
			fmtutil.PrintJSON(mk)

			mv, err := hist.MapKeyValues(mapKey, true)
			logutil.FatalErr(err)
			fmtutil.PrintJSON(mv)
			panic("Z")
		}
	*/
	for mapKeysStr, count := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(mapKeysStr)
		if err != nil {
			return nil, err
		}
		newBinMap := map[string]string{}
		histName := ""
		for k, v := range binMap {
			if k == mapKey {
				histName = v
			} else {
				newBinMap[k] = v
			}
		}
		if len(mapValInclMap) > 0 {
			if _, ok := mapValInclMap[histName]; !ok {
				continue
			}
		}
		subHist, ok := hs.HistogramMap[histName]
		if !ok {
			subHist = NewHistogram(histName)
		}
		subHist.AddMap(newBinMap, count)
		hs.HistogramMap[histName] = subHist
	}
	return hs, nil
}

// ReduceMapKeys returns a new `Histogram` with only the supplied keys present.
func (hist *Histogram) ReduceMapKeys(mapKeysFilter []string) (*Histogram, error) {
	mapKeysFilter = stringsutil.SliceCondenseSpace(mapKeysFilter, true, true)
	filtered := NewHistogram(hist.Name)
	if len(mapKeysFilter) == 0 || len(hist.Bins) == 0 {
		return filtered, nil
	}

	for mapKeysStr, count := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(mapKeysStr)
		if err != nil {
			return nil, err
		}
		newBinMap := map[string]string{}
		for _, filterKey := range mapKeysFilter {
			if filterVal, ok := binMap[filterKey]; ok {
				newBinMap[filterKey] = filterVal
			} else {
				newBinMap[filterKey] = ""
			}
		}
		filtered.AddMap(newBinMap, count)
	}
	return filtered, nil
}

func (hist *Histogram) FlattenMapSingleKey(mapKeyFilter string) (*Histogram, error) {
	filtered := NewHistogram(hist.Name)
	for mapKeyStr, count := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(mapKeyStr)
		if err != nil {
			return nil, err
		}
		if val, ok := binMap[mapKeyFilter]; ok {
			filtered.Add(val, count)
		} else {
			filtered.Add("", count)
		}
	}
	return filtered, nil
}

// MapKeys returns a list of keys using query string keys.
func (hist *Histogram) MapKeys() ([]string, error) {
	keys := map[string]int{}
	for qry := range hist.Bins {
		m, err := maputil.ParseMapStringString(qry)
		if err != nil {
			return []string{}, err
		}
		for k := range m {
			keys[k]++
		}
	}
	return maputil.Keys(keys), nil
}

// MapKeyValues returns a list of keys using query string keys.
func (hist *Histogram) MapKeyValues(key string, dedupe bool) ([]string, error) {
	vals := []string{}
	for qry := range hist.Bins {
		m, err := maputil.ParseMapStringString(qry)
		if err != nil {
			return []string{}, err
		}
		if v, ok := m[key]; ok {
			vals = append(vals, v)
		}
	}
	if dedupe {
		vals = slicesutil.Dedupe(vals)
	}
	return vals, nil
}

func (hist *Histogram) TableSetMap(cfgs []HistogramMapTableConfig) (*table.TableSet, error) {
	if len(cfgs) == 0 {
		return nil, errors.New("error in `Histogram.TableSetMap`: table configs cannot be empty")
	}
	ts := table.NewTableSet(hist.Name)
	for i, cfg := range cfgs {
		if strings.TrimSpace(cfg.SplitKey) == "" {
			tblName := cfg.TableName
			if strings.TrimSpace(tblName) == "" {
				tblName = "Sheet" + strconv.Itoa(i)
			}
			tbl, err := hist.TableMap(cfg.ColumnKeys, cfg.ColNameCount)
			if err != nil {
				return nil, errorsutil.Wrapf(err, "build TableSetMap for (%s)", tblName)
			}

			if slices.Index(ts.Order, tblName) > -1 {
				return nil, errors.New("table name collision")
			}
			tbl.Name = tblName
			ts.Order = append(ts.Order, tblName)
			ts.TableMap[tbl.Name] = tbl
		} else {
			cfg.Inflate()

			hset, err := hist.MapKeySplit(cfg.SplitKey, cfg.SplitValFilterIncl)
			if err != nil {
				return nil, err
			}

			hsetKeys := hset.ItemNames()

			/*
				fmtutil.PrintJSON(hsetKeys)
				fmt.Printf("SPLIT KEY (%s) LEN(%d)\n", cfg.SplitKey, len(hset.HistogramMap))
				panic("HERE")
			*/
			for _, hsetKey := range hsetKeys {
				keyHist, ok := hset.HistogramMap[hsetKey]
				if !ok {
					panic("key not found")
				}
				tbl, err := keyHist.TableMap(cfg.ColumnKeys, cfg.ColNameCount)
				if err != nil {
					return nil, err
				}
				tblName := cfg.TableNamePrefix + hsetKey
				ts.Order = append(ts.Order, tblName)
				tbl.Name = tblName
				ts.TableMap[tblName] = tbl
			}
		}
	}

	return ts, nil
}

type HistogramMapTableSetConfig struct {
	Configs []HistogramMapTableConfig
}

type HistogramMapTableConfig struct {
	TableName          string
	TableNamePrefix    string
	SplitKey           string
	SplitValFilterIncl []string // if present, only include these split values
	ColumnKeys         []string // doesn't include count column
	ColNameCount       string
	splitValFilterMap  map[string]int
	// ColumnNames     []string
}

func (cfg *HistogramMapTableConfig) Inflate() {
	cfg.SplitValFilterIncl = stringsutil.SliceCondenseSpace(cfg.SplitValFilterIncl, true, true)
	cfg.splitValFilterMap = map[string]int{}
	for _, k := range cfg.SplitValFilterIncl {
		cfg.splitValFilterMap[k]++
	}
}

func (cfg *HistogramMapTableConfig) SplitValFilterInclExists(v string) bool {
	if len(cfg.SplitValFilterIncl) != len(cfg.splitValFilterMap) {
		cfg.Inflate()
	}
	if _, ok := cfg.splitValFilterMap[v]; ok {
		return true
	} else {
		return false
	}
}

// TableMap is used to generate a table using map keys.
func (hist *Histogram) TableMap(mapCols []string, colNameBinCount string) (*table.Table, error) {
	if strings.TrimSpace(colNameBinCount) == "" {
		colNameBinCount = "Count"
	}

	// create histogram with minimized aggregate map keys to aggregate exclude non-desired
	// properties from the key for aggregation.
	histSubset := NewHistogram("")
	for binName, binCount := range hist.Bins {
		binMap, err := maputil.ParseMapStringString(binName)
		if err != nil {
			return nil, err
		}
		newBinMap := binMap.Subset(mapCols, false, true, true)
		// newBinMap := mapStringStringSubset(binMap, mapCols, true, false, true)
		// fmtutil.PrintJSON(newBinMap)
		histSubset.AddMap(newBinMap, binCount)
	}

	tbl := table.NewTable(hist.Name)
	tbl.Columns = append(mapCols, colNameBinCount)

	for binName, binCount := range histSubset.Bins {
		binMap, err := maputil.ParseMapStringString(binName)
		if err != nil {
			return nil, err
		}
		binVals := binMap.Gets(true, mapCols)

		tbl.Rows = append(tbl.Rows,
			append(binVals, strconv.Itoa(binCount)),
		)
	}

	tbl.FormatMap = map[int]string{len(tbl.Columns) - 1: "int"}
	return &tbl, nil
}
