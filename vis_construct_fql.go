package main

import (
	"bytes"
	"fmt"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
	"reflect"
	"strings"
)

// select * from orders
// select id, name from orders

type fqlConstructor struct{
	tableName string
	fieldsPartStart string
	fieldsPartEnd string
	fields []string
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
	b.WriteString(v.fieldsPartStart)
	if len(v.fields) > 0 {
		b.WriteString("{")
		b.WriteString(strings.Join(v.fields, ","))
		b.WriteString("}")
	}
	b.WriteString(v.fieldsPartEnd)
	b.WriteString(v.statementTypeEnd)
	return b.String()
}

func (v *fqlConstructor) Enter(in ast.Node) (ast.Node, bool) {
	if reflect.TypeOf(in).String() == "*ast.SelectStmt" {
		v.statementTypeStart = "Map("
	}

	if reflect.TypeOf(in).String() == "*ast.FieldList" {
		if len(in.(*ast.FieldList).Fields) == 1 && len(in.(*ast.FieldList).Fields[0].Text()) == 0 {
			v.fieldsPartStart = "Lambda('x', Get(Var('x'))"
		} else {
			v.fieldsPartStart = "Lambda('x', Let({row: Get(Var('x'))},"
		}
	}

	if reflect.TypeOf(in).String() == "*ast.ColumnName" {
		columnName := in.(*ast.ColumnName).Name.L
		v.fields = append(v.fields, fmt.Sprintf("%s: Select(['data','%s'], Var('row'))", columnName, columnName))
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

	if reflect.TypeOf(in).String() == "*ast.FieldList" {
		if len(in.(*ast.FieldList).Fields) > 1 || len(in.(*ast.FieldList).Fields[0].Text()) > 0 {
			v.fieldsPartEnd = "))"
		}
	}

	return in, true
}