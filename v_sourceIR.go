package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

type collectionIR struct {
	name string
}

func (c *collectionIR) FQLRepr() string {
	return fmt.Sprintf("Collection('%s')", c.name)
}

type indexIR struct {
	name string
}

func (i *indexIR) FQLRepr() string {
	return fmt.Sprintf("Index('%s')", i.name)
}

type sourceIRVisitor struct {
	root fqlIR
}

func (v *sourceIRVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		if len(node.IndexHints) > 0 {
			v.root = &indexIR{node.IndexHints[0].IndexNames[0].L}
		} else {
			v.root = &collectionIR{node.Name.L}
		}
		return in, true
	default:
		return in, false
	}
}

func (v *sourceIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
