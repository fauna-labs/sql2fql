// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

type setFieldIR struct {
	field *fieldIR
	value fqlIR
}

func (f setFieldIR) FQLRepr() string {
	return fmt.Sprintf("%s:%s", f.field.name, f.value.FQLRepr())
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
	case *ast.Assignment:
		next1 := &fqlVisitor{}
		_, _ = node.Expr.Accept(next1)
		v.root.value = next1.root

		next2 := &fieldVisitor{}
		_, _ = node.Column.Accept(next2)
		v.root.field = next2.root
		return in, true
	default:
		return in, false
	}
}

func (v *setFieldVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
