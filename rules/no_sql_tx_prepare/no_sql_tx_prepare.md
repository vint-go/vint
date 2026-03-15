---
title: noSqlTxPrepare
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlTxPrepare`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlTxPrepare:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Tx).Prepare` without a context. The method `(*sql.Tx).Prepare` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Tx).PrepareContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel statement preparation within a transaction, propagate deadlines, and integrate with distributed tracing systems.

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
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    // (*sql.Tx).Prepare does not accept a context
    stmt, err := tx.Prepare("SELECT * FROM users WHERE id = $1")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    defer stmt.Close()
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
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    stmt, err := tx.PrepareContext(ctx, "SELECT * FROM users WHERE id = $1")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    defer stmt.Close()
    tx.Commit()
}
```
