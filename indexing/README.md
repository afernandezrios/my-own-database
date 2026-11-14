# Indexing

## Load test

Baseline performance of the database with no indexes. Every operation goes
through the real pipeline (`engine.Process` → tokenizer → parser → backend), so
the timings include parsing and result formatting.

### Run

From the repository root:

```sh
go -C indexing run .
```

### What it measures

- **INSERT:** 10, 10,000, and 100,000 rows, one SQL statement each.
- **SELECT:** 100,000 point queries (`SELECT * FROM users_N WHERE id = X`) at
  each table size.

Output is a table with total time, average per operation, and ops/sec. `Count`
is the number of operations measured, not result rows. The engine's output is
discarded while measuring so it doesn't flood the terminal, but formatting still
counts toward the time.

The SELECT cost grows with table size because the backend scans every row. That is the cost indexing should remove. After adding an index, expect INSERT to
get slower (index maintenance) and SELECT to get faster.

### Without Index

```
+-----------+--------+--------+-----------------+------------+---------+
| Operation |   Rows |  Count |           Total |     Avg/op | Ops/sec |
+-----------+--------+--------+-----------------+------------+---------+
| INSERT    |     10 |     10 |        18.789µs |    1.878µs |  532226 |
| INSERT    |  10000 |  10000 |     18.847265ms |    1.884µs |  530581 |
| INSERT    | 100000 | 100000 |    185.603485ms |    1.856µs |  538783 |
| SELECT    |     10 | 100000 |    280.324722ms |    2.803µs |  356729 |
| SELECT    |  10000 | 100000 |    7.272202701s |   72.722µs |   13751 |
| SELECT    | 100000 | 100000 | 1m35.927431073s |  959.274µs |    1042 |
+-----------+--------+--------+-----------------+------------+---------+
```

INSERT stays flat (~1.9µs) regardless of size. 
SELECT scales linearly with the scan: 2.7µs → 73µs → 1.01ms.

### With Index

TODO...