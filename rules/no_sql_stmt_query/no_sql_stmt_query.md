---
title: noSqlStmtQuery
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlStmtQuery`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlStmtQuery:
    # no additional options
```

## Details

Disallow calling `(*database/sql.Stmt).Query` without a context. The method `(*sql.Stmt).Query` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.Conn).QueryContext` instead, which accepts a `context.Context` as its first argument.

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
    stmt, err := db.PrepareContext(ctx, "SELECT id, name FROM users WHERE active = $1")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    // (*sql.Stmt).Query does not accept a context
    rows, err := stmt.Query(true)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Println(id, name)
    }
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
    rows, err := conn.QueryContext(ctx, "SELECT id, name FROM users WHERE active = $1", true)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Println(id, name)
    }
}
```
