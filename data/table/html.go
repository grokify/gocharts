package table

import (
	"fmt"
	"html"
	"strings"
)

func (tbl *Table) ToDocuments() []map[string]interface{} {
	docs := []map[string]interface{}{}
	fmtFunc := tbl.FormatterFunc()
	for _, row := range tbl.Records {
		doc := map[string]interface{}{}
		for x, valStr := range row {
			colName := fmt.Sprintf("col%d", x)
			if x < len(tbl.Columns) {
				colNameTry := strings.TrimSpace(tbl.Columns[x])
				if len(colNameTry) > 0 {
					colName = colNameTry
				}
			}
			valFmt, err := fmtFunc(valStr, uint(x))
			if err != nil {
				doc[colName] = valStr
			} else {
				doc[colName] = valFmt
			}
		}
		docs = append(docs, doc)
	}
	return docs
}

// ToHTML converts `*TableData` to HTML.
func (tbl *Table) ToHTML(escapeHTML bool) string {
	tHTML := "<table>"
	tbl.ID = strings.TrimSpace(tbl.ID)
	tbl.Class = strings.TrimSpace(tbl.Class)
	tbl.Style = strings.TrimSpace(tbl.Style)
	attrs := []string{}
	if len(tbl.ID) > 0 {
		attrs = append(attrs, fmt.Sprintf("id=\"%s\"", tbl.ID))
	}
	if len(tbl.Class) > 0 {
		attrs = append(attrs, fmt.Sprintf("class=\"%s\"", tbl.Class))
	}
	if len(tbl.Style) > 0 {
		attrs = append(attrs, fmt.Sprintf("style=\"%s\"", tbl.Style))
	}
	if len(attrs) > 0 {
		tHTML = fmt.Sprintf("<table %s>", strings.Join(attrs, " "))
	}
	if len(tbl.Columns) > 0 {
		if escapeHTML {
			cols := []string{}
			for _, col := range tbl.Columns {
				cols = append(cols, html.EscapeString(col))
			}
			tHTML += "<thead><tr><th>" + strings.Join(cols, "</th><th>") + "</th></tr></thead>"
		} else {
			tHTML += "<thead><tr><th>" + strings.Join(tbl.Columns, "</th><th>") + "</th></tr></thead>"
		}
	}
	if len(tbl.Records) > 0 {
		tHTML += "<tbody>"
		fmtFunc := tbl.FormatterFunc()
		for _, row := range tbl.Records {
			tHTML += "<tr>"
			for x, cell := range row {
				cfmt, err := fmtFunc(cell, uint(x))
				if err != nil {
					if escapeHTML {
						tHTML += "<td>" + html.EscapeString(cell) + "</td>"
					} else {
						tHTML += "<td>" + cell + "</td>"
					}
				} else {
					if escapeHTML {
						tHTML += "<td>" + html.EscapeString(cfmt.(string)) + "</td>"
					} else {
						tHTML += "<td>" + cfmt.(string) + "</td>"
					}
				}
			}
			tHTML += "</tr>"
		}
		tHTML += "</tbody>"
	}
	tHTML += "</table>"
	return tHTML
}
