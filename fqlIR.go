package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

func constructIR(root *ast.StmtNode) fqlIR {
	v := &selectIRVisitor{}
	(*root).Accept(v)
	return v.root
}

type fqlIR interface {
	FQLRepr() string
}

type fqlIRVisitor struct {
	root fqlIR
}

func (v *fqlIRVisitor) Enter(in ast.Node) (res ast.Node, skip bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		next := &selectIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.ColumnName:
		next := &fieldIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.TableName:
		next := &collectionIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.BinaryOperationExpr:
		next := &eqOpIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case ast.ValueExpr:
		next := &valueIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	default:
		res, skip = in, false
	}

	return
}

func (v *fqlIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type valueIR string

func (v valueIR) FQLRepr() string {
	return string(v)
}

type valueIRVisitor struct {
	root valueIR
}

func (v *valueIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case ast.ValueExpr:
		switch value := node.GetValue().(type) {
		case int64:
			v.root = valueIR(fmt.Sprint(value))
		default:
			panic("scalar value not supported")
		}
		return in, true
	default:
		return in, false
	}
}

func (v *valueIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type fieldIR struct {
	name string
}

func (f fieldIR) FQLRepr() string {
	return fmt.Sprintf("Select(['data','%s'], Var('doc'))", f.name)
}

type fieldIRVisitor struct {
	root *fieldIR
}

func (v *fieldIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.ColumnName:
		v.root = &fieldIR{node.Name.L}
		return in, true
	default:
		return in, false
	}
}

func (v *fieldIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type collectionIR struct {
	name string
}

func (c *collectionIR) FQLRepr() string {
	return fmt.Sprintf("Documents(Collection('%s'))", c.name)
}

type collectionIRVisitor struct {
	root *collectionIR
}

func (v *collectionIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		v.root = &collectionIR{node.Name.L}
		return in, true
	default:
		return in, false
	}
}

func (v *collectionIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type eqOpIR struct {
	leftIR  fqlIR
	rightIR fqlIR
}

func (eq eqOpIR) FQLRepr() string {
	return fmt.Sprintf("Equals(%s, %s)", eq.leftIR.FQLRepr(), eq.rightIR.FQLRepr())
}

type eqOpIRVisitor struct {
	root *eqOpIR
}

func (v *eqOpIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.BinaryOperationExpr:
		left := &fqlIRVisitor{}
		right := &fqlIRVisitor{}
		node.L.Accept(left)
		node.R.Accept(right)
		v.root = &eqOpIR{left.root, right.root}
		return in, true
	default:
		return in, false
	}
}

func (v *eqOpIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

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
			filter := &eqOpIRVisitor{}
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
