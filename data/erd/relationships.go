package erd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/er"
)

type Relationships []Relationship

func (rels Relationships) AddToDiagram(d *er.Diagram) *er.Diagram {
	if d == nil {
		d = er.NewDiagram(nil)
	}
	for _, rel := range rels {
		d = d.Relationship(
			rel.LeftEntity, rel.RightEntity,
			rel.LeftRelationship, rel.RightRelationship,
			rel.Identidy, rel.Comment)
	}
	return d
}

func (rels Relationships) DiagramString() string {
	d := rels.AddToDiagram(er.NewDiagram(nil))
	return d.String()
}

func (rels Relationships) MarkdownString(title string) (string, error) {
	b := new(strings.Builder)
	md := markdown.NewMarkdown(b)
	if title != "" {
		md = md.H2(title)
	}
	err := md.CodeBlocks(markdown.SyntaxHighlightMermaid, rels.DiagramString()).
		Build()
	if err != nil {
		return "", err
	} else {
		return b.String(), nil
	}
}

func (rels Relationships) WriteFileMarkdown(filename, title string, perm os.FileMode) error {
	if s, err := rels.MarkdownString(title); err != nil {
		return err
	} else {
		return os.WriteFile(filename, []byte(s), perm)
	}
}

func (rels Relationships) WriteFileHTML(filename, htmlTitle, diagramTitle string, perm os.FileMode) error {
	s := rels.htmlPage(htmlTitle, diagramTitle)
	return os.WriteFile(filename, []byte(s), perm)
}

func (rels Relationships) htmlPage(htmlTitle, diagramTitle string) string {
	md := rels.DiagramString()
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>%s</title>
</head>
<body>
  <h2>%s</h2>
  <div class="mermaid">
%s
  </div>

  <script type="module">
    import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
    mermaid.initialize({ startOnLoad: true });
  </script>
</body>
</html>`
	return fmt.Sprintf(tmpl, htmlTitle, diagramTitle, md)
}
