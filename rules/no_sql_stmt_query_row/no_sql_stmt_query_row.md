---
title: noSqlStmtQueryRow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlStmtQueryRow`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlStmtQueryRow:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Stmt).QueryRow` without a context. The method `(*sql.Stmt).QueryRow` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Conn).QueryRowContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running queries, propagate deadlines, and integrate with distributed tracing systems.

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
    stmt, err := db.PrepareContext(ctx, "SELECT name FROM users WHERE id = $1")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    // (*sql.Stmt).QueryRow does not accept a context
    var name string
    err = stmt.QueryRow(1).Scan(&name)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(name)
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
    conn, err := db.Conn(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    var name string
    err = conn.QueryRowContext(ctx, "SELECT name FROM users WHERE id = $1", 1).Scan(&name)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(name)
}
```
