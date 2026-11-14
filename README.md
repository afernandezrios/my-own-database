# my-own-database

A small SQL database written in Go. Every command goes through this pipeline:

**Engine → Tokenizer → Parser → Backend**

- **Engine:** Receives the SQL command and runs it through the whole pipeline, dispatching the parsed statement to the backend.
- **Tokenizer:** Breaks the raw SQL string into typed tokens (keywords, identifiers, numbers, strings, and symbols).
- **Parser:** Validates the token stream and builds the `Statement` (AST) that the engine executes.
- **Backend:** Keeps the in-memory catalog of tables and performs the create, insert, and select operations.

## Roadmap

The next planned tasks are in [roadmap.md](roadmap.md).
