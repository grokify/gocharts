// tabulator provides helper methods for rendering HTML
// with Tabulator (http://tabulator.info/)
package tabulator

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/grokify/gocharts/v2/data/table"
)

type PageParams struct {
	PageTitle          string
	PageLink           string
	TableDomID         string
	Table              table.Table
	ColumnSet          *ColumnSet
	TabulatorColumnSet *TabulatorColumnSet
	TableJSON          []byte
}

func (pp *PageParams) Inflate() error {
	docs := pp.Table.ToDocuments()
	jdocs, err := json.Marshal(docs)
	if err != nil {
		return err
	}
	pp.TableJSON = jdocs
	return nil
}

func (pp *PageParams) PageLinkHTML() string {
	pp.PageLink = strings.TrimSpace(pp.PageLink)
	if pp.PageLink != "" {
		return html.EscapeString(pp.PageTitle)
	}
	return fmt.Sprintf("<a href=\"%s\">%s</a>",
		pp.PageLink,
		html.EscapeString(pp.PageTitle))
}

func (pp *PageParams) TableJSONBytesOrEmpty() []byte {
	empty := []byte("[]")
	if len(pp.TableJSON) > 0 {
		return pp.TableJSON
	}
	return empty
}

func (pp *PageParams) TabulatorColumnsJSONBytesOrEmpty() []byte {
	if pp.TabulatorColumnSet != nil {
		return pp.TabulatorColumnSet.MustColumnsJSON()
	}
	if pp.ColumnSet == nil || len(pp.ColumnSet.Columns) == 0 {
		// colSet := openapi3.OpTableColumnsDefault(false)
		colset := NewColumnSetSimple(pp.Table.Columns, 1000)
		tcols := colset.Columns.TabulatorColumnSet()
		// tcols := BuildColumnsTabulator(pp.Table.Columns)
		return tcols.MustColumnsJSON()
	}
	tcols := BuildColumnsTabulator(pp.ColumnSet.Columns)
	return tcols.MustColumnsJSON()
}

func (pp *PageParams) WriteFile(filename string) error {
	err := pp.Inflate()
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	WriteTabulatorPage(f, *pp)
	return nil
}
