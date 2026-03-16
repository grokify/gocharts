package chartir

import "strconv"

// ColumnType defines the data type for a column.
type ColumnType string

const (
	ColumnTypeString ColumnType = "string"
	ColumnTypeNumber ColumnType = "number"
)

// Column defines a typed column in a dataset.
type Column struct {
	// Name is the column name/header.
	Name string `json:"name"`

	// Type is the data type for values in this column.
	Type ColumnType `json:"type"`
}

// Dataset represents tabular data for chart consumption.
// The structure is intentionally simple and uniform - always
// typed column definitions plus rows of string values.
type Dataset struct {
	// ID uniquely identifies this dataset for reference by marks.
	ID string `json:"id"`

	// Columns defines the typed column definitions for the data.
	Columns []Column `json:"columns"`

	// Rows contains the data values as strings. Each row is an array
	// of string values corresponding to the columns. The compiler
	// parses values to numbers based on the column type.
	Rows [][]string `json:"rows"`
}

// ColumnIndex returns the index of the column with the given name, or -1 if not found.
func (d *Dataset) ColumnIndex(name string) int {
	for i, col := range d.Columns {
		if col.Name == name {
			return i
		}
	}
	return -1
}

// GetColumn returns the column with the given name, or nil if not found.
func (d *Dataset) GetColumn(name string) *Column {
	idx := d.ColumnIndex(name)
	if idx < 0 {
		return nil
	}
	return &d.Columns[idx]
}

// GetStringValues returns all string values for the given column name.
func (d *Dataset) GetStringValues(colName string) []string {
	idx := d.ColumnIndex(colName)
	if idx < 0 {
		return nil
	}
	values := make([]string, len(d.Rows))
	for i, row := range d.Rows {
		if idx < len(row) {
			values[i] = row[idx]
		}
	}
	return values
}

// GetFloat64Values returns all numeric values for the given column name.
// Non-numeric values are returned as 0.
func (d *Dataset) GetFloat64Values(colName string) []float64 {
	idx := d.ColumnIndex(colName)
	if idx < 0 {
		return nil
	}
	values := make([]float64, len(d.Rows))
	for i, row := range d.Rows {
		if idx < len(row) {
			if v, err := strconv.ParseFloat(row[idx], 64); err == nil {
				values[i] = v
			}
		}
	}
	return values
}
