---
title: noTlsDial
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTlsDial`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTlsDial:
    # no additional options
```

## Details

Disallow calling `crypto/tls.Dial` without a context. The function `tls.Dial` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*crypto/tls.Dialer).DialContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to TLS dial operations enables the caller to cancel connection attempts, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "crypto/tls"
    "fmt"
)

func main() {
    // tls.Dial does not accept a context
    conn, err := tls.Dial("tcp", "example.com:443", &tls.Config{})
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    fmt.Println("Connected")
}
```

### Valid

```golang
package main

import (
    "context"
    "crypto/tls"
    "fmt"
)

func main() {
    ctx := context.Background()
    d := tls.Dialer{
        Config: &tls.Config{},
    }
    conn, err := d.DialContext(ctx, "tcp", "example.com:443")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    fmt.Println("Connected")
}
```
