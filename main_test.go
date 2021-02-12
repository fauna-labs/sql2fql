package main

import "testing"

func TestSimpleSelect(t *testing.T) {
	sql := "select * from a"
	fql := "Map(Paginate(Documents(Collection('a'))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSimpleOptimizedSelect(t *testing.T) {
	sql := "select * from a"
	fql := "Paginate(Match(Index('a_with_values)))"
	assertSQL2FQL(t, sql, fql, true)
}

func TestSelectSingleField(t *testing.T) {
	sql := "select a from b"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))},{a: Select(['data','a'], Var('doc'))})))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectMuiltipleField(t *testing.T) {
	sql := "select a, b from c"
	fql := "Map(Paginate(Documents(Collection('c'))), Lambda('x', Let({doc: Get(Var('x'))},{a: Select(['data','a'], Var('doc')),b: Select(['data','b'], Var('doc'))})))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqual(t *testing.T) {
	sql := "select * from c where a = 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc')), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereGreater(t *testing.T) {
	sql := "select * from c where a > 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, GT(Select(['data','a'], Var('doc')), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereLesser(t *testing.T) {
	sql := "select * from c where a < 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, LT(Select(['data','a'], Var('doc')), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualString(t *testing.T) {
	sql := "select * from c where a = 'hello'"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc')), 'hello'))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualFloatingPoint(t *testing.T) {
	sql := "select * from c where a = 5.0"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc')), 5.0))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualWithLogicalAnd(t *testing.T) {
	sql := "select * from c where a = 5 and b = 6"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, And(Equals(Select(['data','a'], Var('doc')), 5), Equals(Select(['data','b'], Var('doc')), 6)))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualWithLogicalOr(t *testing.T) {
	sql := "select * from c where a = 5 or b = 6"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Or(Equals(Select(['data','a'], Var('doc')), 5), Equals(Select(['data','b'], Var('doc')), 6)))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereGreaterFloatingPoint(t *testing.T) {
	sql := "select * from c where a > 5.0"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, GT(Select(['data','a'], Var('doc')), 5.0))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectWithAnIndex(t *testing.T) {
	sql := "select * from c use index (foo)"
	fql := "Map(Paginate(Match(Index('foo'))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectWithAnIndexAndWhereClause(t *testing.T) {
	sql := "select * from c use index (foo) where a = 5"
	fql := "Map(Paginate(Match(Index('foo'), 5)), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectWithAnIndexAndWhereClauseWithLogicalAnd(t *testing.T) {
	sql := "select * from c use index (foo) where a = 5 and b = 6"
	fql := "Map(Paginate(Match(Index('foo'), 5, 6)), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectWithAnIndexAndWhereClauseWith3LogicalAnd(t *testing.T) {
	sql := "select * from d use index (foo) where a = 5 and b = 6 and c = 7"
	fql := "Map(Paginate(Match(Index('foo'), 5, 6, 7)), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestDeleteWithSingleExactWhereEqual(t *testing.T) {
	sql := "delete from c where a = 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc')), 5))))), Lambda('x', Delete(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func assertSQL2FQL(t *testing.T, sql, fql string, optimize bool) {
	ast, err := parseSql(sql)
	if err != nil {
		t.Error(err)
	}
	actual := ""
	if optimize {
		actual = constructIROptimized(ast).FQLRepr()
	} else {
		actual = constructIR(ast).FQLRepr()
	}
	if actual != fql {
		t.Errorf("\n  actual: %s\nexpected: %s", actual, fql)
	}
}
