package main

import (
	"fmt"

	"github.com/pingcap/parser/ast"
)

func constructIROptimized(root *ast.StmtNode) fqlIR {
	v := &OptimizedFqlVisitor{}
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

type OptimizedFqlVisitor struct {
	optimize bool
	root     fqlIR
}

func (v *OptimizedFqlVisitor) Enter(in ast.Node) (res ast.Node, skip bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		columnsExtractor := &ColumnExtractorVisitor{}
		node.From.Accept(columnsExtractor)
		v.root = &voFqlStatement{}
		return in, false
	default:
		return in, false
	}
}

func (v *OptimizedFqlVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
