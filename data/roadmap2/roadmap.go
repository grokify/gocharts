package roadmap2

import (
	"errors"
	"sort"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// Roadmap is a generic data structure to hold roadmap items.
type Roadmap struct {
	Name         string
	StreamNames  []string
	Columns      table.Columns // Includes Stream Title
	Items        []Item
	ItemCellFunc func(i Item) (colIdx int) // negative means don't use. colIdx is without Stream Title Column
}

// Cell represents a cell of Roadmap Items where defined by a `ColumnName` and `RowName` where the row name
// is the Roadmap stream name.
type Cell struct {
	ColumnName string
	RowName    string
	Items      Items
}

type Cells []Cell

func (c Cells) MaxCellItemsLength() int {
	max := 0
	for _, cell := range c {
		if len(cell.Items) > max {
			max = len(cell.Items)
		}
	}
	return max
}

func (r *Roadmap) ItemsByStream(streamName string) []Item {
	itms := []Item{}
	for _, item := range r.Items {
		if item.StreamName == streamName {
			itms = append(itms, item)
		}
	}
	return itms
}

func (r *Roadmap) streamsMap() map[string]int {
	streamsMap := map[string]int{}
	for _, streamNaem := range r.StreamNames {
		streamsMap[streamNaem] = 0
	}
	return streamsMap
}

func (r *Roadmap) UnknownStreams() []string {
	streamsMap := r.streamsMap()
	unknownStreamsMap := map[string]int{}
	for _, itm := range r.Items {
		if len(itm.StreamName) == 0 {
			continue
		}
		if _, ok := streamsMap[itm.StreamName]; !ok {
			unknownStreamsMap[itm.StreamName] += 1
		}
	}
	if len(unknownStreamsMap) == 0 {
		return []string{}
	}
	unknownStreamNames := []string{}
	for k := range unknownStreamsMap {
		unknownStreamNames = append(unknownStreamNames, k)
	}
	sort.Strings(unknownStreamNames)
	return unknownStreamNames
}

// Stream represents the horizontal swimlane on a roadmap chart.
type Stream struct {
	Name  string
	Cells Cells
}

// Streams returns a slice of `Stream`.
func (r *Roadmap) Streams(includeUnknown, includeUnassigned bool) ([]Stream, error) {
	streams := []Stream{}
	if r.ItemCellFunc == nil {
		return streams, errors.New("itemCellFunc cannot be nil")
	}
	for _, knownStreamName := range r.StreamNames {
		stream := Stream{
			Name:  knownStreamName,
			Cells: []Cell{},
		}
		streamItems := r.ItemsByStream(knownStreamName)
		streamCellsArray := make([]Cell, len(r.Columns)-1)
		for _, itm := range streamItems {
			colIdx := r.ItemCellFunc(itm)
			if colIdx < 0 {
				continue
			}
			if colIdx >= len(r.Columns)-1 {
				return streams, errorsutil.ErrIndexOutOfRange(colIdx, len(r.Columns)-1)
			}
			cell := streamCellsArray[colIdx]
			cell.Items = append(cell.Items, itm)
			streamCellsArray[colIdx] = cell
		}
		stream.Cells = append(stream.Cells, streamCellsArray...)
		for i, cell := range stream.Cells {
			cell.RowName = knownStreamName
			if i+1 < len(r.Columns) {
				cell.ColumnName = r.Columns[i+1]
			}
			stream.Cells[i] = cell
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

// Table returns a `table` representing a roadmap slide. It can be used to generate CSV or XLSX files.
func (r *Roadmap) Table(includeUnknown, includeUnassigned bool) (table.Table, error) {
	tbl := table.NewTable(r.Name)
	tbl.Columns = r.Columns

	streams, err := r.Streams(includeUnknown, includeUnassigned)
	if err != nil {
		return tbl, err
	}
	for _, stream := range streams {
		maxCellLength := stream.Cells.MaxCellItemsLength()
		if maxCellLength == 0 {
			continue
		}

		rows := slicesutil.MakeMatrix2D[string](maxCellLength, len(tbl.Columns))
		for x, cell := range stream.Cells {
			for _, item := range cell.Items {
				y, err := stringsutil.Matrix2DColRowIndex(rows, uint(x+1), "")
				if err != nil {
					return tbl, err
				}
				if y < 0 {
					panic("empty index not found")
				}
				rows[y][x+1] = item.Name
			}
		}
		for y := range rows {
			rows[y][0] = stream.Name
		}
		tbl.Rows = append(tbl.Rows, rows...)
	}

	return tbl, nil
}
