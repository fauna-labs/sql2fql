package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

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
