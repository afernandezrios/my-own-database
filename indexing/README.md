# Indexing

## Load test

End-to-end performance of the database. Every operation goes through the real
pipeline (`engine.Process` → tokenizer → parser → backend), so the timings
include parsing and result formatting. Tables with an `id` column carry an
in-memory B-tree index (`backend/btree_index.go`) that is created automatically
with the table and maintained on every INSERT.

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

### Without Index (baseline)

Recorded before the `id` index was wired in. The SELECT cost grows linearly with
the table because the backend scans every row — the cost indexing removes.

```
+-----------+--------+--------+-----------------+------------+---------+
| Operation |   Rows |  Count |           Total |     Avg/op | Ops/sec |
+-----------+--------+--------+-----------------+------------+---------+
| INSERT    |     10 |     10 |         17.83µs |    1.783µs |  560852 |
| INSERT    |  10000 |  10000 |     20.507697ms |     2.05µs |  487622 |
| INSERT    | 100000 | 100000 |    180.665711ms |    1.806µs |  553508 |
| SELECT    |     10 | 100000 |     282.16296ms |    2.821µs |  354405 |
| SELECT    |  10000 | 100000 |    8.558399199s |   85.583µs |   11684 |
| SELECT    | 100000 | 100000 | 1m47.661464308s | 1.076614ms |     929 |
+-----------+--------+--------+-----------------+------------+---------+
```

INSERT stays flat (~1.8–2.1µs) regardless of size.
SELECT scales linearly with the scan: 2.8µs → 86µs → 1.08ms.

### With Index

Same machine, same program, with the `id` index maintained on INSERT and used
for `WHERE id = N` lookups.

```
+-----------+--------+--------+--------------+----------+---------+
| Operation |   Rows |  Count |        Total |   Avg/op | Ops/sec |
+-----------+--------+--------+--------------+----------+---------+
| INSERT    |     10 |     10 |     30.781µs |  3.078µs |  324876 |
| INSERT    |  10000 |  10000 | 23.747499ms |  2.374µs |  421097 |
| INSERT    | 100000 | 100000 | 204.429983ms |  2.044µs |  489165 |
| SELECT    |     10 | 100000 |  254.39789ms |  2.543µs |  393085 |
| SELECT    |  10000 | 100000 | 265.340951ms |  2.653µs |  376874 |
| SELECT    | 100000 | 100000 | 276.675604ms |  2.766µs |  361434 |
+-----------+--------+--------+--------------+----------+---------+
```

INSERT is slightly slower than the baseline — index maintenance adds a key parse
and a tree insert per row — but stays flat with table size, as an O(log n)
insert predicts. SELECT is now flat at ~2.5–2.8µs regardless of table size; the
end-to-end floor is dominated by tokenizing, parsing, and formatting, not the
lookup. At 100,000 rows a point query went from ~1.08ms to ~2.8µs (~390×).