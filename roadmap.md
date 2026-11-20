# Roadmap

## Core

1. **BTree index + benchmark:** Add an `Index` interface to the backend, maintain it on insert, and use it for equality filters; benchmark scan vs index at 10 vs 10k rows to see read speed vs write amplification.
2. **Persistence via WAL + checkpoint:** Append every successful statement to a WAL, replay it on startup, and snapshot every N operations so startup stays bounded; measure fsync on vs off and verify recovery after a `kill -9`.
3. **Transactions over the WAL:** `BEGIN` buffers operations, `COMMIT` appends a commit record and applies them, `ROLLBACK` discards them; replay ignores uncommitted transactions, and index changes roll back with the rows.

## Stretch

- **In-memory B-tree:** Replace the hash index with a B-tree to learn splits, merges, and fanout.
- **Isolation demo:** Expose the engine to goroutine clients in tests, reproduce lost updates, then compare a global lock with row locks.
- **Page storage + buffer pool:** Move rows into slotted pages with pinning and an eviction policy — the start of a second iteration.

Note: the global `catalog` in `src/engine/engine.go` needs to become a `DB`/session object that owns the WAL and transaction state.
