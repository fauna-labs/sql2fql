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

	switch s.source.(type) {
	case *collectionIR:
		if s.filter != nil {
			filter := "Filter(Documents(%s), Lambda('x', Let({doc: Get(Var('x'))}, %s)))"
			sb.WriteString(fmt.Sprintf(filter, s.source.FQLRepr(), s.filter.FQLRepr()))
		} else {
			sb.WriteString(fmt.Sprintf("Documents(%s)", s.source.FQLRepr()))
		}
	case *indexIR:
		if s.filter != nil {
			b := s.filter.(*binaryOperatorIR)
			if b.op != EQ {
				panic("indexes only works with equality operators")
			}

			sb.WriteString(fmt.Sprintf("Match(%s, %s)", s.source.FQLRepr(), b.rightIR.FQLRepr()))
		} else {
			sb.WriteString(fmt.Sprintf("Match(%s)", s.source.FQLRepr()))
		}
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

type selectVisitor struct {
	root *selectIR
}

func (v *selectVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		v.root = &selectIR{}

		source := &sourceVisitor{}
		node.From.Accept(source)
		v.root.source = source.root

		for _, fNode := range node.Fields.Fields {
			if fNode.Expr != nil {
				field := &fieldVisitor{}
				fNode.Expr.Accept(field)
				v.root.fields = append(v.root.fields, field.root)
			}
		}

		if node.Where != nil {
			filter := &binaryOperatorVisitor{}
			node.Where.Accept(filter)
			v.root.filter = filter.root
		}
		return in, true
	default:
		return in, false
	}
}

func (v *selectVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
