package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

func constructIROptimized(root *ast.StmtNode) fqlIR {
	v := &voFqlIROptimized{}
	(*root).Accept(v)
	return v.root
}

type voFqlIR interface {
	FQLRepr() string
}

type voFqlStatement struct {
}

func (c *voFqlStatement) FQLRepr() string {
	return fmt.Sprintf("Paginate(Match(Index('a_with_values)))")
}

type voFqlIROptimized struct {
	optimize bool
	root     fqlIR
}

func (v *voFqlIROptimized) Enter(in ast.Node) (res ast.Node, skip bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		columnsExtractor := &voColumnExtractor{}
		node.From.Accept(columnsExtractor)
		fmt.Println(columnsExtractor.colNames)
		v.root = &voFqlStatement{}
		return in, false
	default:
		return in, false
	}
}

func (v *voFqlIROptimized) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
