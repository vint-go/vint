---
title: noNetDialTimeout
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetDialTimeout`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetDialTimeout:
    # no additional options
```

## Details

Disallow calling `net.DialTimeout` without a context. The function `net.DialTimeout` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.Dialer).DialContext` with `(*net.Dialer).Timeout` instead, which accepts a `context.Context` as its first argument and allows setting a timeout via the `Dialer` struct.

Passing `context.Context` to dial operations enables the caller to cancel connection attempts, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "net"
    "time"
)

func main() {
    // net.DialTimeout does not accept a context
    conn, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```

### Valid

```golang
package main

import (
    "context"
    "net"
    "time"
)

func main() {
    ctx := context.Background()
    d := net.Dialer{
        Timeout: 5 * time.Second,
    }
    conn, err := d.DialContext(ctx, "tcp", "example.com:80")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```
