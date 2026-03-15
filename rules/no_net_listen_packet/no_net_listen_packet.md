---
title: noNetListenPacket
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetListenPacket`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetListenPacket:
    # no additional options
```

## Details

Disallow calling `net.ListenPacket` without a context. The function `net.ListenPacket` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.ListenConfig).ListenPacket` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to network operations enables the caller to cancel operations, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import "net"

func main() {
    // net.ListenPacket does not accept a context
    conn, err := net.ListenPacket("udp", ":8080")
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
    lc := net.ListenConfig{}
    conn, err := lc.ListenPacket(ctx, "udp", ":8080")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
}
```
