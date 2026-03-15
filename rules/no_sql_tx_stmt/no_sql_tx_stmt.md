---
title: noSqlTxStmt
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlTxStmt`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlTxStmt:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Tx).Stmt` without a context. The method `(*sql.Tx).Stmt` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Tx).StmtContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel statement association within a transaction, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "context"
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "connstring")
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()
    stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    // (*sql.Tx).Stmt does not accept a context
    txStmt := tx.Stmt(stmt)
    _ = txStmt
    tx.Commit()
}
```

### Valid

```golang
package main

import (
    "context"
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "connstring")
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()
    stmt, err := db.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    txStmt := tx.StmtContext(ctx, stmt)
    _ = txStmt
    tx.Commit()
}
```
