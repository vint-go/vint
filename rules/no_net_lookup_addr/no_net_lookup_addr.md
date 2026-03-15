---
title: noNetLookupAddr
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetLookupAddr`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetLookupAddr:
    # no additional options
```

## Details

Disallow calling `net.LookupAddr` without a context. The function `net.LookupAddr` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.Resolver).LookupAddr` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to DNS lookup operations enables the caller to cancel long-running lookups, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "net"
)

func main() {
    // net.LookupAddr does not accept a context
    names, err := net.LookupAddr("127.0.0.1")
    if err != nil {
        panic(err)
    }
    fmt.Println(names)
}
```

### Valid

```golang
package main

import (
    "context"
    "fmt"
    "net"
)

func main() {
    ctx := context.Background()
    r := net.Resolver{}
    names, err := r.LookupAddr(ctx, "127.0.0.1")
    if err != nil {
        panic(err)
    }
    fmt.Println(names)
}
```
