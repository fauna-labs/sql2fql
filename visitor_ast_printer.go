package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
)

func printAst(rootNode *ast.StmtNode) {
	v := &astPrinterVisitor{depth: 0}
	(*rootNode).Accept(v)
	return
}

type astPrinterVisitor struct {
	depth int
}

func (v *astPrinterVisitor) Enter(in ast.Node) (ast.Node, bool) {
	v.depth = v.depth + 1
	prefix := strings.Repeat("   ", v.depth)
	fmt.Printf("%v %v\n", prefix, reflect.TypeOf(in))

	switch node := in.(type) {
	case *ast.IndexOption:
		fmt.Printf("%v -> %#v\n", prefix, node)
	}

	return in, false
}

func (v *astPrinterVisitor) Leave(in ast.Node) (ast.Node, bool) {
	v.depth = v.depth - 1
	return in, true
}
