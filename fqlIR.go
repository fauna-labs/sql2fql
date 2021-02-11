package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

type fqlIR interface {
	FQLRepr() string
}

type fieldIR struct {
	name string
}

func (f fieldIR) FQLRepr() string {
	return fmt.Sprintf("Select(['data','%s'], Var('doc'))", f.name)
}

type collectionIR string

func (c collectionIR) FQLRepr() string {
	return fmt.Sprintf("Documents(Collection('%s'))", c)
}

type selectIR struct {
	source fqlIR
	fields []*fieldIR
	filter fqlIR
}

type eqOpIR struct {
	leftIR  fqlIR
	rightIR fqlIR
}

func (eq eqOpIR) FQLRepr() string {
	return fmt.Sprintf("Equals(%s, %s)", eq.leftIR.FQLRepr(), eq.rightIR.FQLRepr())
}

func (s *selectIR) FQLRepr() string {
	var sb strings.Builder

	sb.WriteString("Map(Paginate(")

	if s.filter != nil {
		filter := "Filter(%s, Lambda('x', Let({doc: Get(Var('x'))}, %s)))"
		sb.WriteString(fmt.Sprintf(filter, s.source.FQLRepr(), s.filter.FQLRepr()))
	} else {
		sb.WriteString(s.source.FQLRepr())
	}

	sb.WriteString("), ")

	if len(s.fields) == 0 {
		sb.WriteString("Lambda('x', Get(Var('x'))")
	} else {
		sb.WriteString("Lambda('x', Let({doc: Get(Var('x'))},{")

		for i, f := range s.fields {
			sb.WriteString(fmt.Sprintf("%s: %s", f.name, f.FQLRepr()))
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
		field := &fieldIR{node.Name.L}
		sel := f.root.(*selectIR)
		sel.fields = append(sel.fields, field)

	case *ast.TableName:
		sel := f.root.(*selectIR)
		sel.source = collectionIR(node.Name.L)

	case *ast.BinaryOperationExpr:
		sel := f.root.(*selectIR)
		sel.filter = &eqOpIR{}
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
