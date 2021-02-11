package main

import (
	"github.com/pingcap/parser/ast"
)

type voColumnExtractor struct {
	colNames []string
}

func (v *voColumnExtractor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.TableName:
		v.colNames = append(v.colNames, node.Name.O)
		return in, true
	default:
		return in, false
	}
}

func (v *voColumnExtractor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
