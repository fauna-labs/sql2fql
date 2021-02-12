package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
)

type binaryOperatorIR struct {
	op      operation
	leftIR  fqlIR
	rightIR fqlIR
}

func (bin binaryOperatorIR) FQLRepr() string {
	switch bin.op {
	case EQ:
		return fmt.Sprintf("Equals(%s, %s)", bin.leftIR.FQLRepr(), bin.rightIR.FQLRepr())
	case GT:
		return fmt.Sprintf("GT(%s, %s)", bin.leftIR.FQLRepr(), bin.rightIR.FQLRepr())
	case GTE:
		return fmt.Sprintf("GTE(%s, %s)", bin.leftIR.FQLRepr(), bin.rightIR.FQLRepr())
	case LT:
		return fmt.Sprintf("LT(%s, %s)", bin.leftIR.FQLRepr(), bin.rightIR.FQLRepr())
	case LTE:
		return fmt.Sprintf("LTE(%s, %s)", bin.leftIR.FQLRepr(), bin.rightIR.FQLRepr())
	}
	panic("Unsupported binary operation type")
}

type binaryOperatorVisitor struct {
	root *binaryOperatorIR
}

func (v *binaryOperatorVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.BinaryOperationExpr:
		left := &fqlVisitor{}
		right := &fqlVisitor{}
		node.L.Accept(left)
		node.R.Accept(right)
		op := getOperation(in.(*ast.BinaryOperationExpr).Op)
		v.root = &binaryOperatorIR{op: op, leftIR: left.root, rightIR: right.root}
		return in, true
	default:
		return in, false
	}
}

func (v *binaryOperatorVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type operation int

func getOperation(op opcode.Op) operation {
	switch op {
	case opcode.EQ:
		return EQ
	case opcode.GT:
		return GT
	case opcode.GE:
		return GTE
	case opcode.LT:
		return LT
	case opcode.LE:
		return LTE
	}
	return ERR
}

const (
	ERR = -1
	EQ  = 0
	GT  = 1
	GTE = 2
	LT  = 3
	LTE = 4
)