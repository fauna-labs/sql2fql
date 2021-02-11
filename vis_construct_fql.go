package main

import (
	"bytes"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"reflect"
)

// select * from orders
// select id, name from orders

type fqlConstructor struct{
	tableName string
	fieldsPart string
	statementTypeStart string
	statementTypeEnd string
}

func constructAst(rootNode *ast.StmtNode) string {
	v := &fqlConstructor { }
	(*rootNode).Accept(v)

	return joinFqlParts(v)
}

func joinFqlParts(v *fqlConstructor) string {
	var b bytes.Buffer
	b.WriteString(v.statementTypeStart)
	b.WriteString(v.tableName)
	b.WriteString(",")
	b.WriteString(v.fieldsPart)
	b.WriteString(v.statementTypeEnd)
	return b.String()
}

func (v *fqlConstructor) Enter(in ast.Node) (ast.Node, bool) {
	if reflect.TypeOf(in).String() == "*ast.SelectStmt" {
		v.statementTypeStart = "Map("
	}
	if reflect.TypeOf(in).String() == "*ast.FieldList" {
		if len(in.(*ast.FieldList).Fields) == 1 {
			v.fieldsPart = "Lambda('x', Get(Var('x'))"
		} else {
			// implement Let logic statement
			v.fieldsPart = ""
		}
	}
	if reflect.TypeOf(in).String() == "*ast.TableName" {
		v.tableName = "Paginate(Documents(Collection('" + in.(*ast.TableName).Name.L + "')))"
	}

	return in, false
}

func (v *fqlConstructor) Leave(in ast.Node) (ast.Node, bool) {
	if reflect.TypeOf(in).String() == "*ast.SelectStmt" {
		v.statementTypeEnd = ")"
	}

	return in, true
}