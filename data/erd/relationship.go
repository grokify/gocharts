package erd

import (
	"github.com/nao1215/markdown/mermaid/er"
)

type Relationship struct {
	LeftEntity        er.Entity
	RightEntity       er.Entity
	LeftRelationship  er.Relationship
	RightRelationship er.Relationship
	Identidy          er.Identify
	Comment           string
}
