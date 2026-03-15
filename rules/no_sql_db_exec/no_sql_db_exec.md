---
title: noSqlDbExec
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlDbExec`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlDbExec:
    # no additional options
```

## Details

Disallow calling `(*database/sql.DB).Exec` without a context. The method `(*sql.DB).Exec` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.DB).ExecContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running queries, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "connstring")
    if err != nil {
        log.Fatal(err)
    }
    // (*sql.DB).Exec does not accept a context
    _, err = db.Exec("INSERT INTO users (name) VALUES ($1)", "Alice")
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
    _, err = db.ExecContext(ctx, "INSERT INTO users (name) VALUES ($1)", "Alice")
    if err != nil {
        log.Fatal(err)
    }
}
```
