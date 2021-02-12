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
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))},{a: Select(['data','a'], Var('doc'), null)})))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleField(t *testing.T) {
	sql := "update b SET a=5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:5}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateMultipleField(t *testing.T) {
	sql := "update b SET a=5, b=4"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:5,b:4}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithSingleExactWhereEqual(t *testing.T) {
	sql := "update b SET a=5 where a=7"
	fql := "Map(Paginate(Filter(Documents(Collection('b')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc'), null), 7))))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:5}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithPlusExpression(t *testing.T) {
	sql := "update b SET a=a+5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:Sum([Select(['data','a'], Var('doc'), null), 5])}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithMinusExpression(t *testing.T) {
	sql := "update b SET a=a-5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:Subtract(Select(['data','a'], Var('doc'), null), 5)}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithMultiplyExpression(t *testing.T) {
	sql := "update b SET a=a*5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:Multiply(Select(['data','a'], Var('doc'), null), 5)}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithDivideExpression(t *testing.T) {
	sql := "update b SET a=a/5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:Divide(Select(['data','a'], Var('doc'), null), 5)}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestUpdateSingleFieldWithModuloExpression(t *testing.T) {
	sql := "update b SET a=a%5"
	fql := "Map(Paginate(Documents(Collection('b'))), Lambda('x', Let({doc: Get(Var('x'))}, Update(Var('x'), {data:{a:Modulo(Select(['data','a'], Var('doc'), null), 5)}}))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectMuiltipleField(t *testing.T) {
	sql := "select a, b from c"
	fql := "Map(Paginate(Documents(Collection('c'))), Lambda('x', Let({doc: Get(Var('x'))},{a: Select(['data','a'], Var('doc'), null),b: Select(['data','b'], Var('doc'), null)})))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqual(t *testing.T) {
	sql := "select * from c where a = 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc'), null), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereGreater(t *testing.T) {
	sql := "select * from c where a > 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, GT(Select(['data','a'], Var('doc'), null), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereLesser(t *testing.T) {
	sql := "select * from c where a < 5"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, LT(Select(['data','a'], Var('doc'), null), 5))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualString(t *testing.T) {
	sql := "select * from c where a = 'hello'"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc'), null), 'hello'))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualFloatingPoint(t *testing.T) {
	sql := "select * from c where a = 5.0"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc'), null), 5.0))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualWithLogicalAnd(t *testing.T) {
	sql := "select * from c where a = 5 and b = 6"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, And(Equals(Select(['data','a'], Var('doc'), null), 5), Equals(Select(['data','b'], Var('doc'), null), 6)))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereEqualWithLogicalOr(t *testing.T) {
	sql := "select * from c where a = 5 or b = 6"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Or(Equals(Select(['data','a'], Var('doc'), null), 5), Equals(Select(['data','b'], Var('doc'), null), 6)))))), Lambda('x', Get(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestSelectSingleFieldWithSingleExactWhereGreaterFloatingPoint(t *testing.T) {
	sql := "select * from c where a > 5.0"
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, GT(Select(['data','a'], Var('doc'), null), 5.0))))), Lambda('x', Get(Var('x'))))"
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
	fql := "Map(Paginate(Filter(Documents(Collection('c')), Lambda('x', Let({doc: Get(Var('x'))}, Equals(Select(['data','a'], Var('doc'), null), 5))))), Lambda('x', Delete(Var('x'))))"
	assertSQL2FQL(t, sql, fql, false)
}

func TestInsertQuery(t *testing.T) {
	sql := "insert into a (b, c) values ('foo', 'bar')"
	fql := "Create(Collection('a'), { data: { b: 'foo', c: 'bar' }})"
	assertSQL2FQL(t, sql, fql, false)
}

func TestCreateCollection(t *testing.T) {
	sql := "create table a"
	fql := "CreateCollection({ name: 'a' })"
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
