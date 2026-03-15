---
title: noNetDial
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetDial`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetDial:
    # no additional options
```

## Details

Disallow calling `net.Dial` without a context. The function `net.Dial` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.Dialer).DialContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to dial operations enables the caller to cancel connection attempts, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import "net"

func main() {
    // net.Dial does not accept a context
    conn, err := net.Dial("tcp", "example.com:80")
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
)

func main() {
    ctx := context.Background()
    d := net.Dialer{}
    conn, err := d.DialContext(ctx, "tcp", "example.com:80")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```
