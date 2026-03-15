---
title: noNetLookupCname
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNetLookupCname`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNetLookupCname:
    # no additional options
```

## Details

Disallow calling `net.LookupCNAME` without a context. The function `net.LookupCNAME` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net.Resolver).LookupCNAME` instead, which accepts a `context.Context` as its first argument.

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
    // net.LookupCNAME does not accept a context
    cname, err := net.LookupCNAME("www.example.com")
    if err != nil {
        panic(err)
    }
    fmt.Println(cname)
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
    cname, err := r.LookupCNAME(ctx, "www.example.com")
    if err != nil {
        panic(err)
    }
    fmt.Println(cname)
}
```
