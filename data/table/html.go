package table

import (
	"fmt"
	"html"
	"strings"
)

func ToDocuments(tbl *Table) []map[string]interface{} {
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
func ToHTML(tbl *Table, domID string, escapeHTML bool) string {
	tHTML := "<table>"
	domID = strings.TrimSpace(domID)
	if len(domID) > 0 {
		tHTML = fmt.Sprintf("<table id=\"%s\">", domID)
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
