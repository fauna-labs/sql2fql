package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
)

type selectIR struct {
	statement statementType
	source    fqlIR
	fields    []*fieldIR
	setFields []*setFieldIR
	filter    fqlIR
}

type statementType int

func (s *selectIR) FQLRepr() string {
	var sb strings.Builder

	sb.WriteString("Map(Paginate(")

	switch s.source.(type) {
	case *collectionIR:
		if s.filter != nil {
			filter := "Filter(Documents(%s), Lambda('x', Let({doc: Get(Var('x'))}, %s)))"
			sb.WriteString(fmt.Sprintf(filter, s.source.FQLRepr(), s.filter.FQLRepr()))
		} else {
			sb.WriteString(fmt.Sprintf("Documents(%s)", s.source.FQLRepr()))
		}
	case *indexIR:
		if s.filter != nil {
			b := s.filter.(*binaryOperatorIR)
			sb.WriteString(fmt.Sprintf("Match(%s, %s)", s.source.FQLRepr(), strings.Join(indexValues(b, nil), ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("Match(%s)", s.source.FQLRepr()))
		}
	}

	sb.WriteString("), ")

	switch s.statement {
	case SELECT:
		fallthrough
	case DELETE:
		writeSelectDeleteFields(s, &sb)
	case UPDATE:
		if len(s.setFields) == 0 {
			panic("UPDATE without SET")
		} else {
			sb.WriteString(fmt.Sprintf("Lambda('x', Update({"))
			// write fields.
			// Update(Var('x'), {data:{a:5}})

			sb.WriteString("}))")
		}
	}

	sb.WriteString(")")
	return sb.String()
}

func indexValues(b *binaryOperatorIR, res []string) []string {
	switch b.op {
	case EQ:
		res = append(res, b.rightIR.FQLRepr())
	case LOGIC_AND:
		res = indexValues(b.leftIR.(*binaryOperatorIR), res)
		res = indexValues(b.rightIR.(*binaryOperatorIR), res)
	default:
		panic("indexes only works with equality operators")
	}
	return res
}

type selectVisitor struct {
	root *selectIR
}

func (v *selectVisitor) Enter(in ast.Node) (ast.Node, bool) {
	switch node := in.(type) {
	case *ast.SelectStmt:
		handleSelect(v, node.From, node.Where, node.Fields.Fields, SELECT)
		return in, true
	case *ast.DeleteStmt:
		handleSelect(v, node.TableRefs, node.Where, []*ast.SelectField{}, DELETE)
		return in, true
	case *ast.UpdateStmt:
		handleUpdate(v, node.TableRefs, node.Where, node.List)
		return in, true
	default:
		return in, false
	}
}

func (v *selectVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func handleSelect(v *selectVisitor, from *ast.TableRefsClause, where ast.ExprNode, fields []*ast.SelectField, t statementType) {
	v.root = &selectIR{statement: t}
	source := &sourceVisitor{}
	from.Accept(source)
	v.root.source = source.root

	for _, fNode := range fields {
		if fNode.Expr != nil {
			field := &fieldVisitor{}
			fNode.Expr.Accept(field)
			v.root.fields = append(v.root.fields, field.root)
		}
	}

	if where != nil {
		filter := &binaryOperatorVisitor{}
		where.Accept(filter)
		v.root.filter = filter.root
	}
}

func handleUpdate(v *selectVisitor, from *ast.TableRefsClause, where ast.ExprNode, setFields []*ast.Assignment) {
	v.root = &selectIR{statement: UPDATE}
	source := &sourceVisitor{}
	from.Accept(source)
	v.root.source = source.root

	for _, fNode := range setFields {
		setField := &setFieldVisitor{}
		fNode.Accept(setField)
		fmt.Println(setField)
		v.root.setFields = append(v.root.setFields, setField.root)

	}

	if where != nil {
		filter := &binaryOperatorVisitor{}
		where.Accept(filter)
		v.root.filter = filter.root
	}
}

func writeSelectDeleteFields(s *selectIR, sb *strings.Builder) {
	action := ""
	switch s.statement {
	case SELECT:
		action = "Get"
	case DELETE:
		action = "Delete"
	}
	if len(s.fields) == 0 {
		sb.WriteString(fmt.Sprintf("Lambda('x', %s(Var('x')))", action))
	} else {
		sb.WriteString(fmt.Sprintf("Lambda('x', Let({doc: %s(Var('x'))},{", action))

		for i, f := range s.fields {
			sb.WriteString(fmt.Sprintf("%s: %s", f.name, f.FQLRepr()))
			if i < len(s.fields)-1 {
				sb.WriteString(",")
			}
		}

		sb.WriteString("}))")
	}
}

const (
	SELECT = 0
	DELETE = 1
	UPDATE = 2
)
