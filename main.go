package main

import (
	"os"
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"

)

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: colx 'SQL statement'")
		return
	}
	sql := os.Args[1]
	astNode, err := parse(sql)

	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	printAst(astNode)
	fmt.Printf("Columns: %v\n", extractColumns(astNode))
}

