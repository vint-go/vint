---
title: noSqlTxQueryRow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlTxQueryRow`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlTxQueryRow:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Tx).QueryRow` without a context. The method `(*sql.Tx).QueryRow` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Tx).QueryRowContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running queries within a transaction, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "context"
    "database/sql"
    "fmt"
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
    // (*sql.Tx).QueryRow does not accept a context
    var name string
    err = tx.QueryRow("SELECT name FROM users WHERE id = $1", 1).Scan(&name)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    fmt.Println(name)
    tx.Commit()
}
```

### Valid

```golang
package main

import (
    "context"
    "database/sql"
    "fmt"
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
    var name string
    err = tx.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }
    fmt.Println(name)
    tx.Commit()
}
```
