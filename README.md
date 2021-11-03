This repository contains unofficial patterns, sample code, or tools to help
developers build more effectively with [Fauna][fauna]. All [Fauna
Labs][fauna-labs] repositories are provided “as-is” and without support. By
using this repository or its contents, you agree that this repository may never
be officially supported and moved to the [Fauna
organization][fauna-organization].

---

# Fauna SQL to FQL utility

If you are familiar with SQL and want to ease into FQL (Fauna Query Language), this repository will help you execute simple SQL statements on your Fauna database by translating them into FQL Queries. While it is not possible to provide an exhaustive translation between SQL and FQL, this is a useful tool for you to learn and get started with FQL. As you gain more experience with the FQL, the ease and power of its syntax should become evident. Complex queries that are difficult, or even impossible, in SQL can be composed very easily in FQL.

## Requirements

* The [Go][golang] programming language +1.15
* The [fauna-shell][fauna-shell] utility (optional)

## Using this repository

To use SQL to FQL utility, clone this repository to your local machine and build
the program with the `go build` command. Once built, a command line tool called
`sql2fql` will become available in the repo's root directory. If the `--key`
command line argument is provided, `sql2fql` will attempt to execute the
generated FQL query using [fauna-shell][fauna-shell] utility. Additional command
line arguments can be found with the `./sql2fql -h` command.

## Supported SQL statements

Bellow are some examples of supported SQL statements.

### CREATE TABLE

```bash
./sql2fql --sql "create table users"

 SQL  create table users

 FQL  CreateCollection({
          name: 'users'
      })
```

### CREATE INDEX

```bash
./sql2fql --sql "create index user_by_name on users (name)"

 SQL  create index user_by_name on users (name);

 FQL  CreateIndex({
          name: 'user_by_name',
          source: Collection('users'),
          unique: false,
          terms: [{
              field: ['data', 'name']
          }]
      })
```

### SELECT

```bash
./sql2fql --sql "select * from users"

 SQL  select * from users

 FQL  Map(Paginate(Documents(Collection('users'))), Lambda('x', Get(Var('x'))))
```

```bash
./sql2fql --sql "select name, age from users"

 SQL  select name, age from users

 FQL  Map(Paginate(Documents(Collection('users'))), Lambda('x', Let({
          doc: Get(Var('x'))
      }, {
          name: Select(['data', 'name'], Var('doc'), null),
          age: Select(['data', 'age'], Var('doc'), null)
      })))
```

```bash
./sql2fql --sql "select name, age from users where name = 'bob'"

 SQL  select name, age from users where name = 'bob'

 FQL  Map(Paginate(Filter(Documents(Collection('users')), Lambda('x', Let({
          doc: Get(Var('x'))
      }, Equals(Select(['data', 'name'], Var('doc'), null), 'bob'))))), Lambda('x', Let({
          doc: Get(Var('x'))
      }, {
          name: Select(['data', 'name'], Var('doc'), null),
          age: Select(['data', 'age'], Var('doc'), null)
      })))
```

```bash
./sql2fql --sql "select name, age from users use index (user_by_name) where name = 'bob'"

 SQL  select name, age from users use index (user_by_name) where name = 'bob'

 FQL  Map(Paginate(Match(Index('user_by_name'), 'bob')), Lambda('x', Let({
          doc: Get(Var('x'))
      }, {
          name: Select(['data', 'name'], Var('doc'), null),
          age: Select(['data', 'age'], Var('doc'), null)
      })))
```

### INSERT

```bash
./sql2fql --sql "insert into users (name, age) values ('bob', 42)"

 SQL  insert into users (name, age) values ('bob', 42)

 FQL  Create(Collection('users'), {
          data: {
              name: 'bob',
              age: 42
          }
      })
```

### UPDATE

```bash
./sql2fql --sql "update users set age = 43 where name = 'bob'"

 SQL  update users set age = 43 where name = 'bob'

 FQL  Map(Paginate(Filter(Documents(Collection('users')), Lambda('x', Let({
          doc: Get(Var('x'))
      }, Equals(Select(['data', 'name'], Var('doc'), null), 'bob'))))), Lambda('x', Let({
          doc: Get(Var('x'))
      }, Update(Var('x'), {
          data: {
              age: 43
          }
      }))))
```

```bash
./sql2fql --sql "update users use index (user_by_name) set age = 43 where name = 'bob'"

 SQL  update users use index (user_by_name) set age = 43 where name = 'bob'

 FQL  Map(Paginate(Match(Index('user_by_name'), 'bob')), Lambda('x', Let({
          doc: Get(Var('x'))
      }, Update(Var('x'), {
          data: {
              age: 43
          }
      }))))
```

### DELETE

```bash
./sql2fql --sql "delete from users where name = 'bob'"

 SQL  delete from users where name = 'bob'

 FQL  Map(Paginate(Filter(Documents(Collection('users')), Lambda('x', Let({
          doc: Get(Var('x'))
      }, Equals(Select(['data', 'name'], Var('doc'), null), 'bob'))))), Lambda('x', Delete(Var('x'))))
```

```bash
./sql2fql --sql="delete from users use index (user_by_name) where name = 'bob'"

 SQL  delete from users use index (user_by_name) where name = 'bob'

 FQL  Map(Paginate(Match(Index('user_by_name'), 'bob')), Lambda('x', Delete(Var('x'))))
```

---

Copyright Fauna, Inc. or its affiliates. All rights reserved.
SPDX-License-Identifier: MIT-0

[fauna]: https://fauna.com
[fauna-labs]: https://github.com/fauna-labs
[fauna-organization]: https://github.com/fauna
[fauna-shell]: https://github.com/fauna/fauna-shell
[golang]: https://golang.org/
