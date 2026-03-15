---
title: noSqlDbPrepare
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlDbPrepare`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlDbPrepare:
    # no additional options
```

## Details

Disallow calling `(*database/sql.DB).Prepare` without a context. The method `(*sql.DB).Prepare` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.DB).PrepareContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel statement preparation, propagate deadlines, and integrate with distributed tracing systems.

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
    // (*sql.DB).Prepare does not accept a context
    stmt, err := db.Prepare("SELECT * FROM users WHERE id = $1")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
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
}
```
