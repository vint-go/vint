---
title: noSqlStmtExec
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlStmtExec`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlStmtExec:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Stmt).Exec` without a context. The method `(*sql.Stmt).Exec` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Conn).ExecContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running queries, propagate deadlines, and integrate with distributed tracing systems.

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
    stmt, err := db.PrepareContext(ctx, "INSERT INTO users (name) VALUES ($1)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    // (*sql.Stmt).Exec does not accept a context
    _, err = stmt.Exec("Alice")
    if err != nil {
        log.Fatal(err)
    }
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
    conn, err := db.Conn(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    _, err = conn.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
    if err != nil {
        log.Fatal(err)
    }
}
```
