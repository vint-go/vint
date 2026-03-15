---
title: noNetLookupNs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetLookupNs`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetLookupNs:
    # no additional options
```

## Details

Disallow calling `net.LookupNS` without a context. The function `net.LookupNS` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.Resolver).LookupNS` instead, which accepts a `context.Context` as its first argument.

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
    // net.LookupNS does not accept a context
    nss, err := net.LookupNS("example.com")
    if err != nil {
        panic(err)
    }
    for _, ns := range nss {
        fmt.Println(ns.Host)
    }
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
    nss, err := r.LookupNS(ctx, "example.com")
    if err != nil {
        panic(err)
    }
    for _, ns := range nss {
        fmt.Println(ns.Host)
    }
}
```
