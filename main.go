package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"github.com/spf13/cobra"
	"github.com/jedib0t/go-pretty/v6/table"
)

const art = `
             .__   ________   _____      .__   
  ___________|  |  \_____  \_/ ____\_____|  |  
 /  ___/ ____/  |   /  ____/\   __\/ ____/  |  
 \___ < <_|  |  |__/       \ |  | < <_|  |  |__
/____  >__   |____/\_______ \|__|  \__   |____/
     \/   |__|             \/         |__|     `

var rootCommand = &cobra.Command{
	Use:   "sql2fql",
	Short: art,
	Run:   run,
}

const shellCommand = "fauna"

var sql string
var optimize bool
var key string
var colors bool
var tables bool

func main() {
	rootCommand.Flags().StringVarP(&sql, "sql", "s", "", "the SQL shellCommand")
	rootCommand.MarkFlagRequired("sql")
	rootCommand.Flags().BoolVarP(&optimize, "optimize", "o", false, "whether to optimize queries using indexes")
	rootCommand.Flags().StringVarP(&key, "key", "k", "", "the key to use to run the query")
	rootCommand.Flags().BoolVarP(&colors, "color", "c", true, "whether to color output")
	rootCommand.Flags().BoolVarP(&tables, "tables", "b", false, "whether to put output into tables")

	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, _ []string) {
	t := table.NewWriter()
	if !tables {
		t.SetStyle(table.Style{
			Name: "NoBorders",
			Box: table.BoxStyle{
				BottomLeft:       " ",
				BottomRight:      " ",
				BottomSeparator:  " ",
				Left:             " ",
				LeftSeparator:    " ",
				MiddleHorizontal: " ",
				MiddleSeparator:  " ",
				MiddleVertical:   " ",
				PaddingLeft:      " ",
				PaddingRight:     " ",
				Right:            " ",
				RightSeparator:   " ",
				TopLeft:          " ",
				TopRight:         " ",
				TopSeparator:     " ",
				UnfinishedRow:    " ",
			},
		})
	}

	if colors {
		t.AppendRow(table.Row{"SQL", color.CyanString(sql)})
	} else {
		t.AppendRow(table.Row{"SQL", sql})
	}
	fql := transpileSqlToFql(sql)
	t.AppendSeparator()
	if colors {
		t.AppendRow(table.Row{"FQL", color.MagentaString(fql)})
	} else {
		t.AppendRow(table.Row{"FQL", fql})
	}
	if key != "" {
		if !shellInstalled() {
			panic(fmt.Sprintf("fauna shell isn't installed or configured correctly"))
		}
		out := executeFql(fql, key)
		t.AppendSeparator()
		if colors {
			t.AppendRow(table.Row{"Output", color.GreenString(out)})
		} else {
			t.AppendRow(table.Row{"Output", out})
		}
	}
	fmt.Println()
	fmt.Println(t.Render())
}

func transpileSqlToFql(sql string) string {
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
	var out bytes.Buffer
	cmd := exec.Command(shellCommand, "eval", fmt.Sprintf("--secret=%s", key), "--format=shell", fql)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(fmt.Sprintf("error executing fql: %s", err))
	}
	return out.String()
}
