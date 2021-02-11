package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

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
