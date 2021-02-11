package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

type fqlIR interface {
	FQLRepr() string
}

type fieldIR string

func (f fieldIR) FQLRepr() string { return string(f) }

type collectionIR string

func (c collectionIR) FQLRepr() string {
	return fmt.Sprintf("Collection('%s')", c)
}

type selectIR struct {
	source fqlIR
	fields []fieldIR
}

func (s *selectIR) FQLRepr() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Map(Paginate(Documents(%s)),", s.source.FQLRepr()))

	if len(s.fields) == 0 {
		sb.WriteString("Lambda('x', Get(Var('x'))")
	} else {
		sb.WriteString("Lambda('x', Let({row: Get(Var('x'))},{")

		for i, f := range s.fields {
			sb.WriteString(fmt.Sprintf("%s: Select(['data','%s'], Var('row'))", f, f))
			if i < len(s.fields)-1 {
				sb.WriteString(",")
			}
		}

		sb.WriteString("}))")
	}

	sb.WriteString(")")
	return sb.String()
}

type fqlIRVisitor struct {
	root fqlIR
}

func (f *fqlIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		f.root = &selectIR{}

	case *ast.ColumnName:
		sel := f.root.(*selectIR)
		sel.fields = append(sel.fields, fieldIR(node.Name.L))

	case *ast.TableName:
		sel := f.root.(*selectIR)
		sel.source = collectionIR(node.Name.L)
	}

	return in, false
}

func (v *fqlIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func constructIR(root *ast.StmtNode) fqlIR {
	v := &fqlIRVisitor{}
	(*root).Accept(v)
	return v.root
}
