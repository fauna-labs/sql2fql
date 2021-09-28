// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
)

type valueIR string

func (v valueIR) FQLRepr() string {
	return string(v)
}

type valueVisitor struct {
	root valueIR
}

func (v *valueVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case ast.ValueExpr:
		value := node.GetValue()
		switch value.(type) {
		case int64:
			v.root = valueIR(fmt.Sprint(value))
		case string:
			v.root = valueIR(fmt.Sprintf("'%s'", value))
		case *test_driver.MyDecimal:
			v.root = valueIR(fmt.Sprint(value))
		default:
			panic("scalar value not supported")
		}
		return in, true
	default:
		return in, false
	}
}

func (v *valueVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
