package main

import (
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
)

func extractColumns(rootNode *ast.StmtNode) []string {
	v := &columnExtractor{}
	(*rootNode).Accept(v)
	return v.colNames
}

type columnExtractor struct{
	colNames []string
}

func (v *columnExtractor) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *columnExtractor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}