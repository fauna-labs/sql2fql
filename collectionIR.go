package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

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
