package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/parser/test_driver"
)

// select * from orders
// select id, name from orders

type fqlConstructor struct {
	tableName          string
	fieldsPartStart    string
	fieldsPartEnd      string
	fields             []string
	statementTypeStart string
	statementTypeEnd   string
}

func constructAst(rootNode *ast.StmtNode) string {
	v := &fqlConstructor{}
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
	switch node := in.(type) {
	case *ast.SelectStmt:
		v.statementTypeStart = "Map("
	case *ast.FieldList:
		if len(node.Fields) == 1 && len(node.Fields[0].Text()) == 0 {
			v.fieldsPartStart = "Lambda('x', Get(Var('x'))"
		} else {
			v.fieldsPartStart = "Lambda('x', Let({row: Get(Var('x'))},"
		}
	case *ast.ColumnName:
		columnName := node.Name.L
		v.fields = append(v.fields, fmt.Sprintf("%s: Select(['data','%s'], Var('row'))", columnName, columnName))
	case *ast.TableName:
		v.tableName = "Paginate(Documents(Collection('" + node.Name.L + "')))"
	}
	return in, false
}

func (v *fqlConstructor) Leave(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		v.statementTypeEnd = ")"
	case *ast.FieldList:
		if len(node.Fields) > 1 || len(node.Fields[0].Text()) > 0 {
			v.fieldsPartEnd = "))"
		}
	}
	return in, true
}
