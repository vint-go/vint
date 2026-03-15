---
title: noNetListen
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetListen`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetListen:
    # no additional options
```

## Details

Disallow calling `net.Listen` without a context. The function `net.Listen` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.ListenConfig).Listen` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to network operations enables the caller to cancel long-running listeners, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import "net"

func main() {
    // net.Listen does not accept a context
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer ln.Close()
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
    ln, err := lc.Listen(ctx, "tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer ln.Close()
}
```
