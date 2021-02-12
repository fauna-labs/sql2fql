package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

type fieldIR struct {
	name string
}

func (f fieldIR) FQLRepr() string {
	return fmt.Sprintf("Select(['data','%s'], Var('doc'))", f.name)
}

type fieldVisitor struct {
	optimize bool
	root     *fieldIR
}

func (v *fieldVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.ColumnName:
		v.root = &fieldIR{node.Name.L}
		return in, true
	default:
		return in, false
	}
}

func (v *fieldVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
