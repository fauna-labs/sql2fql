package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

type setFieldIR struct {
	name  fqlIR
	value fqlIR
}

func (f setFieldIR) FQLRepr() string {
	return fmt.Sprintf("{data:{ %s:%s}}", f.name.FQLRepr(), f.value.FQLRepr())
}

type setFieldVisitor struct {
	optimize bool
	root     *setFieldIR
}

func (v *setFieldVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if v.root == nil {
		v.root = &setFieldIR{}
	}

	switch node := in.(type) {
	case ast.ValueExpr:
		next := &valueVisitor{}
		_, _ = node.Accept(next)
		v.root.value = next.root
		return in, true

	case *ast.ColumnName:
		next := &fieldVisitor{}
		_, _ = node.Accept(next)
		v.root.name = next.root

		return in, true
	default:
		return in, false
	}
}

func (v *setFieldVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
