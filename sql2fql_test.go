package main

import "testing"

func TestSimpleSelect(t *testing.T) {
	sql := "select * from a"
	fql := "Map(Paginate(Documents(Collection('a'))),Lambda('x', Get(Var('x')))"
	assertSQL2FQL(t, sql, fql)
}

func TestSelectSingleField(t *testing.T) {
	sql := "select a from b"
	fql := "Map(Paginate(Documents(Collection('b'))),Lambda('x', Let({row: Get(Var('x'))},{a: Select(['data','a'], Var('row'))})))"
	assertSQL2FQL(t, sql, fql)
}

func TestSelectMuiltipleField(t *testing.T) {
	sql := "select a, b from c"
	fql := "Map(Paginate(Documents(Collection('c'))),Lambda('x', Let({row: Get(Var('x'))},{a: Select(['data','a'], Var('row')),b: Select(['data','b'], Var('row'))})))"
	assertSQL2FQL(t, sql, fql)
}

func TestSelectSingleFieldWithSingleExactWhere(t *testing.T) {
	sql := "select a from c where a = 5"
	fql := "Map(Paginate(Match(Index('c_by_a'), 5)), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql)
}

func assertSQL2FQL(t *testing.T, sql, fql string) {
	ast, err := parse(sql)
	if err != nil {
		t.Error(err)
	}

	actual := constructIR(ast).FQLRepr()
	if actual != fql {
		t.Errorf("\n  actual: %s\nexpected: %s", actual, fql)
	}
}