package main

import (
	"fmt"
	"os"
	"strconv"

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
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . 'SQL statement' true/false")
		return
	}

	sql := os.Args[1]
	optimize := false

	if len(os.Args) > 2 {
		optimizeStr := os.Args[2]
		opt, err := strconv.ParseBool(optimizeStr)
		if err != nil {
			fmt.Printf("bool parse error: %v\n", err.Error())
			return
		}
		optimize = opt
	}

	astNode, err := parse(sql)

	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}

	//printAst(astNode)
	//fmt.Printf("Columns: %v\n", extractColumns(astNode))
	var ir fqlIR
	if optimize {
		ir = constructIR(astNode)
	} else {
		ir = constructIROptimized(astNode)
	}
	fmt.Println(ir.FQLRepr())
}
