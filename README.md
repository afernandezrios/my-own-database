# my-own-database

A small SQL database written in Go. Every command goes through this pipeline:

**Engine → Tokenizer → Parser → Backend**

- **Engine:** Receives the SQL command and runs it through the whole pipeline, dispatching the parsed statement to the backend.
- **Tokenizer:** Breaks the raw SQL string into typed tokens (keywords, identifiers, numbers, strings, and symbols).
- **Parser:** Validates the token stream and builds the `Statement` (AST) that the engine executes.
- **Backend:** Keeps the in-memory catalog of tables and performs the create, insert, and select operations. Tables with an `id` column also get an in-memory B-tree index.

## Indexing

A table gets a B-tree index when it is created with an `id` column. The backend maintains the index on every insert and uses it for `WHERE id = ...` equality filters. All other filters fall back to a sequential scan.

The [indexing load test](indexing/README.md) compares indexed and unindexed performance.

## Roadmap

The next planned tasks are in [roadmap.md](roadmap.md).
