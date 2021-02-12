package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

type createCollectionIR struct {
	name string
}

func (c *createCollectionIR) FQLRepr() string {
	return fmt.Sprintf("CreateCollection({ name: '%s' })", c.name)
}

type createIndexIR struct {
	name   string
	source string
	terms  []string
	unique bool
}

func (c *createIndexIR) FQLRepr() string {
	var sb strings.Builder
	sb.WriteString("CreateIndex({ name: '")
	sb.WriteString(c.name)
	sb.WriteString("', source: Collection('")
	sb.WriteString(c.source)
	sb.WriteString("'), unique: ")

	if c.unique {
		sb.WriteString("true")
	} else {
		sb.WriteString("false")
	}

	sb.WriteString(", terms: [")

	for i, term := range c.terms {
		sb.WriteString("{field: ['data', '")
		sb.WriteString(term)
		sb.WriteString("']}")
		if i < len(c.terms)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("]})")
	return sb.String()
}

type createVisitor struct {
	root fqlIR
}

func (v *createVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		v.root = &createCollectionIR{node.Table.Name.O}
		return in, true

	case *ast.CreateIndexStmt:
		create := &createIndexIR{}
		create.name = node.IndexName
		create.source = node.Table.Name.O
		create.unique = node.KeyType == ast.IndexKeyTypeUnique

		for _, col := range node.IndexPartSpecifications {
			create.terms = append(create.terms, col.Column.Name.O)
		}

		v.root = create
		return in, true

	default:
		return in, false
	}
}

func (v *createVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
