package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

type createIR struct {
	name string
}

func (c *createIR) FQLRepr() string {
	return fmt.Sprintf("CreateCollection({ name: '%s' })", c.name)
}

type createVisitor struct {
	root *createIR
}

func (v *createVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.CreateTableStmt:
		v.root = &createIR{node.Table.Name.O}
		return in, true
	default:
		return in, false
	}
}

func (v *createVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
