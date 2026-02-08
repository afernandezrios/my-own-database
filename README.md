# my-own-database

---

### The Architecture

Before writing code, visualize the pipeline. Every database command goes through these phases:

1. **REPL:** The interface where the user types.
2. **Lexer:** Breaks the text string into "Tokens" (e.g., `SELECT`, `*`, `FROM`).
3. **Parser:** Arranges Tokens into an Abstract Syntax Tree (AST)—a structure the code understands.
4. **Backend/Engine:** The logic that actually reads or writes the data based on the AST.

---

### The REPL (Read-Eval-Print Loop)

**Goal:** 

Create a CLI that runs, accepts input, and prints it back until the user types `exit`.

**Task:** 

Use `bufio.Scanner` in a `for` loop to capture `stdin`.

**Acceptance Criteria:**

    When you run `go run main.go`, you see a prompt: `db > `.
    Typing `hello` and hitting Enter prints `Unrecognized command "hello"` and immediately shows the prompt `db > ` again.
    Typing `.exit` (or just `exit`, your choice) terminates the program immediately with code 0.

### The Tokenizer (Lexer)

**Goal:** 

Turn a raw string like `"INSERT INTO users VALUES (1, 'Tim')"` into a slice of categorized tokens.

**Task:** 

Create a `Tokenize(input string) ([]Token, error)` function.

**Acceptance Criteria:**

*Note: Since the REPL isn't connected to this yet, verify this with a `lexer_test.go` unit test.*

1. **Input:** `"select id from users"`
2. **Output:** A slice of length 4.
   - `Token{Type: SELECT, Value: "select"}`
   - `Token{Type: IDENT, Value: "id"}`
   - `Token{Type: FROM, Value: "from"}`
   - `Token{Type: IDENT, Value: "users"}`

3. **Edge Case:** Input `"select *"` should recognize `*` as a symbol, not an identifier.

### The Parser (The AST)

**Goal:** 

Validate that the tokens make grammatical sense and organize them into a Struct.

**Task:** 

Create a `Parse(tokens []Token) (*Statement, error)` function.

**Acceptance Criteria:**

*Note: Verify with `parser_test.go*`*

1. **Valid Select:**
   * **Input:** Tokens for `"SELECT name FROM users"`
   * **Output:** Returns a `Statement` struct where `Statement.Type` is "SELECT", `Statement.Fields` contains `["name"]`, and `Statement.TableName` is `"users"`.

2. **Valid Insert:**
   * **Input:** Tokens for `"INSERT INTO users VALUES (1, 'Tim')"`
   * **Output:** `Statement` struct where `Statement.Values` contains `["1", "Tim"]`.

3. **Invalid Syntax:**
   * **Input:** Tokens for `"SELECT FROM users"` (Missing column name)
   * **Output:** Returns a non-nil `error` message: `"Syntax error: expected identifier after SELECT"`.

### The In-Memory Backend

**Goal:** 

Create the Go structs to hold the data in RAM.

**Task:** 

Create `Table` and `Catalog` structs and methods to manipulate them directly (bypassing SQL for now).

**Acceptance Criteria:**

*Note: *Verify with `backend_test.go`*

1. **Creation:** Call `catalog.CreateTable("users", []string{"id", "name"})`. Verify the map `catalog.Tables` now contains a key "users".
2. **Insertion:** Call `table.InsertRow([]string{"1", "Tim"})`. Verify `len(table.Rows)` is 1.
3. **Data Integrity:** Call `table.InsertRow([]string{"1"})` (missing 'name'). The function should return an error because the row length doesn't match the column length.

### The Execution Engine (Integration)

**Goal:** 

Connect the REPL -> Lexer -> Parser -> Backend.

**Task:** 

In your main loop, pass the parsed AST to an `Execute()` function.
* **### How it should work (Acceptance Criteria):**
* *Now you can test manually in the CLI (go run main.go)*

1. **Create:**
   * Type: `CREATE TABLE users (id, name)`
   * Response: `Table 'users' created.`

2. **Insert:**
   * Type: `INSERT INTO users VALUES (1, 'Tim')`
   * Response: `1 row inserted.`

3. **Select:**
   * Type: `SELECT * FROM users`
   * Response:
    ```text
    id   | name
    -----+------
    1    | Tim
    ```

4. **Error Propagation:**
* Type: `INSERT INTO non_existent_table ...`
* Response: `Error: Table 'non_existent_table' does not exist.`

### Filtering (The WHERE Clause)

**Goal:** 

Support `SELECT * FROM users WHERE id = 1`.

**Task:** 

Update Parser to handle `WHERE` and Engine to filter rows.

**Acceptance Criteria:**

1. **Setup:** Insert two rows: `(1, 'Tim')` and `(2, 'Sarah')`.
2. **Match:**
   * Type: `SELECT * FROM users WHERE id = 1`
   * Response: Only prints the row for Tim.

3. **No Match:**
   * Type: `SELECT * FROM users WHERE id = 99`
   * Response: Prints the headers, but 0 rows.

4. **Bad Column:**
   * Type: `SELECT * FROM users WHERE age = 30` (Assume 'age' column doesn't exist)
   * Response: `Error: Column 'age' does not exist in table 'users'.`

---

### To go further:

**Persistence:** Use `encoding/gob` or `encoding/json` to save your Table maps to a file on disk when the program 
closes, and load them when it starts.

**Indexing:** Implementing a simple B-Tree to make lookups faster.