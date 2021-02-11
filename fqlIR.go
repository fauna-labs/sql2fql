package main

import (
	"github.com/pingcap/parser/ast"
)

func constructIR(root *ast.StmtNode) fqlIR {
	v := &selectIRVisitor{}
	(*root).Accept(v)
	return v.root
}

type fqlIR interface {
	FQLRepr() string
}

type fqlIRVisitor struct {
	root fqlIR
}

func (v *fqlIRVisitor) Enter(in ast.Node) (res ast.Node, skip bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		next := &selectIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.ColumnName:
		next := &fieldIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.TableName:
		next := &collectionIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.BinaryOperationExpr:
		next := &eqOpIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case ast.ValueExpr:
		next := &valueIRVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	default:
		res, skip = in, false
	}

	return
}

func (v *fqlIRVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
