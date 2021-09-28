// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package main

import (
	"github.com/pingcap/parser/ast"
)

type ColumnExtractorVisitor struct {
	colNames []string
}

func (v *ColumnExtractorVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		v.colNames = append(v.colNames, node.Name.O)
		return in, true
	default:
		return in, false
	}
}

func (v *ColumnExtractorVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
