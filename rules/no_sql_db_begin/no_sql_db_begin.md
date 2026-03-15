---
title: noSqlDbBegin
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlDbBegin`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlDbBegin:
    # no additional options
```

## Details

Disallow calling `(*database/sql.DB).Begin` without a context. The method `(*sql.DB).Begin` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.DB).BeginTx` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel long-running transactions, propagate deadlines, and integrate with distributed tracing systems.

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
    // (*sql.DB).Begin does not accept a context
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    defer tx.Rollback()
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
    defer tx.Rollback()
}
```
