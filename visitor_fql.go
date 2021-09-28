// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package main

import (
	"github.com/pingcap/parser/ast"
)

func constructIR(root *ast.StmtNode) fqlIR {
	v := &fqlVisitor{}
	(*root).Accept(v)
	return v.root
}

type fqlIR interface {
	FQLRepr() string
}

type fqlVisitor struct {
	optimize bool
	root     fqlIR
}

func (v *fqlVisitor) Enter(in ast.Node) (res ast.Node, skip bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		next := &selectVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.InsertStmt:
		next := &insertVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.CreateTableStmt:
		next := &createVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.CreateIndexStmt:
		next := &createVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.DeleteStmt:
		next := &selectVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.UpdateStmt:
		next := &selectVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.ColumnName:
		next := &fieldVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.TableName:
		next := &sourceVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case *ast.BinaryOperationExpr:
		next := &binaryOperatorVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	case ast.ValueExpr:
		next := &valueVisitor{}
		res, skip = node.Accept(next)
		v.root = next.root

	default:
		res, skip = in, false
	}

	return
}

func (v *fqlVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
