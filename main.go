package main

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/spf13/cobra"
	"os"
)

const art = `
             .__   ________   _____      .__   
  ___________|  |  \_____  \_/ ____\_____|  |  
 /  ___/ ____/  |   /  ____/\   __\/ ____/  |  
 \___ < <_|  |  |__/       \ |  | < <_|  |  |__
/____  >__   |____/\_______ \|__|  \__   |____/
     \/   |__|             \/         |__|     `

var rootCommand = &cobra.Command{
	Use:	"sql2fql",
	Short:	art,
	Run: transpileSqlToFql,
}

var sql string
var optimize bool

func main() {
	rootCommand.Flags().StringVarP(&sql, "sql", "s", "", "the SQL command")
	rootCommand.MarkFlagRequired("sql")
	rootCommand.Flags().BoolVarP(&optimize, "optimize", "o", false, "whether to use indexes")
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func transpileSqlToFql (cmd *cobra.Command, _ []string) {
	node, err := parseSql(sql)
	if err != nil {
		panic(fmt.Sprintf("error parsing sql: %s", err.Error()))
	}
	ir := constructIR(node)
	fmt.Println(ir.FQLRepr())
}

func parseSql(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	nodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}