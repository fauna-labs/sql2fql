package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

type selectIR struct {
	source fqlIR
	fields []*fieldIR
	filter fqlIR
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
		sb.WriteString("Lambda('x', Get(Var('x')))")
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

type selectIRVisitor struct {
	root *selectIR
}

func (v *selectIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		v.root = &selectIR{}

		source := &collectionIRVisitor{}
		node.From.Accept(source)
		v.root.source = source.root

		for _, fNode := range node.Fields.Fields {
			if fNode.Expr != nil {
				field := &fieldIRVisitor{}
				fNode.Expr.Accept(field)
				v.root.fields = append(v.root.fields, field.root)
			}
		}

		if node.Where != nil {
			filter := &binOpIRVisitor{}
			node.Where.Accept(filter)
			v.root.filter = filter.root
		}
		return in, true
	default:
		return in, false
	}
}

func (v *selectIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
