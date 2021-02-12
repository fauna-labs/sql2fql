package main

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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
	Run: run,
}

const shellCommand = "fauna"

var sql string
var optimize bool
var key string

func main() {
	rootCommand.Flags().StringVarP(&sql, "sql", "s", "", "the SQL shellCommand")
	rootCommand.MarkFlagRequired("sql")
	rootCommand.Flags().BoolVarP(&optimize, "optimize", "o", false, "whether to use indexes")
	rootCommand.Flags().StringVarP(&key, "key", "k", "", "the key to use to run the query")
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run (cmd *cobra.Command, _ []string) {
	fql := transpileSqlToFql(sql)
	fmt.Println(fql)
	if key != "" {
		if !shellInstalled() {
			panic(fmt.Sprintf("fauna shell isn't installed or configured correctly"))
		}
		out := executeFql(fql, key)
		fmt.Println(out)
	}
}

func transpileSqlToFql (sql string) string {
	node, err := parseSql(sql)
	if err != nil {
		panic(fmt.Sprintf("error parsing sql: %s", err.Error()))
	}
	var ir fqlIR
	if optimize {
		ir = constructIROptimized(node)
	} else {
		ir = constructIR(node)
	}
	return ir.FQLRepr()
}

func parseSql(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	nodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return &nodes[0], nil
}

func shellInstalled() bool {
	cmd := exec.Command("command", "-v", shellCommand)
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func executeFql(fql string, key string) string {
	cmd := exec.Command(shellCommand, "eval", fmt.Sprintf("--secret=\"%s\"", key), "--format=shell", fmt.Sprintf("\"%s\"", fql))
	fmt.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("error executing fql: %s", err))
	}
	out, _ := cmd.Output()
	return string(out)
}