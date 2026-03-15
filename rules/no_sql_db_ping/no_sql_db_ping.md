---
title: noSqlDbPing
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSqlDbPing`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSqlDbPing:
    # no additional options
```

## Details

Disallow calling `(*database/sql.DB).Ping` without a context. The method `(*sql.DB).Ping` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*database/sql.DB).PingContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to database operations enables the caller to cancel connectivity checks, propagate deadlines, and integrate with distributed tracing systems.

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
    // (*sql.DB).Ping does not accept a context
    if err := db.Ping(); err != nil {
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
    if err := db.PingContext(ctx); err != nil {
        log.Fatal(err)
    }
}
```
