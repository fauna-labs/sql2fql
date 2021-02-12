package main

import (
	"strings"

	"github.com/pingcap/parser/ast"
)

type insertIR struct {
	source *collectionIR
	fields []string
	values []fqlIR
}

func (i *insertIR) FQLRepr() string {
	var sb strings.Builder
	sb.WriteString("Create(")
	sb.WriteString(i.source.FQLRepr())
	sb.WriteString(", { data: { ")

	for idx, field := range i.fields {
		value := i.values[idx]
		sb.WriteString(field)
		sb.WriteString(": ")
		sb.WriteString(value.FQLRepr())
		if idx < len(i.fields)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(" }})")
	return sb.String()
}

type insertVisitor struct {
	root *insertIR
}

func (v *insertVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.InsertStmt:
		v.root = &insertIR{}

		source := &sourceVisitor{}
		node.Table.Accept(source)
		if collIR, ok := source.root.(*collectionIR); ok {
			v.root.source = collIR
		} else {
			panic("can only call insert on a collection")
		}

		for _, col := range node.Columns {
			v.root.fields = append(v.root.fields, col.Name.O)
		}

		for _, lst := range node.Lists[0] {
			value := &fqlVisitor{}
			lst.Accept(value)
			v.root.values = append(v.root.values, value.root)
		}

		return in, true
	default:
		return in, false
	}
}

func (v *insertVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
