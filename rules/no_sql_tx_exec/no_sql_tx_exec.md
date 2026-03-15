---
title: noSqlTxExec
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlTxExec`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlTxExec:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Tx).Exec` without a context. The method `(*sql.Tx).Exec` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Tx).ExecContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running queries within a transaction, propagate deadlines, and integrate with distributed tracing systems.

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
    // (*sql.Tx).Exec does not accept a context
    _, err = tx.Exec("INSERT INTO users (name) VALUES ($1)", "Alice")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
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
    _, err = tx.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    tx.Commit()
}
```
