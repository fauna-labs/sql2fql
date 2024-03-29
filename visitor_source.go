// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

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

type sourceVisitor struct {
	root fqlIR
}

func (v *sourceVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		if len(node.IndexHints) > 0 {
			v.root = &indexIR{node.IndexHints[0].IndexNames[0].O}
		} else {
			v.root = &collectionIR{node.Name.O}
		}
		return in, true
	default:
		return in, false
	}
}

func (v *sourceVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
