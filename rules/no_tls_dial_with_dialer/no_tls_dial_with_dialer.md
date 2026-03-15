---
title: noTlsDialWithDialer
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTlsDialWithDialer`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTlsDialWithDialer:
    # no additional options
```

## Details

Disallow calling `crypto/tls.DialWithDialer` without a context. The function `tls.DialWithDialer` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*crypto/tls.Dialer).DialContext` with `NetDialer` instead, which accepts a `context.Context` as its first argument and allows specifying a custom net dialer.

Passing `context.Context` to TLS dial operations enables the caller to cancel connection attempts, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "crypto/tls"
    "fmt"
    "net"
    "time"
)

func main() {
    dialer := &net.Dialer{
        Timeout: 5 * time.Second,
    }
    // tls.DialWithDialer does not accept a context
    conn, err := tls.DialWithDialer(dialer, "tcp", "example.com:443", &tls.Config{})
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
    "net"
    "time"
)

func main() {
    ctx := context.Background()
    d := tls.Dialer{
        NetDialer: &net.Dialer{
            Timeout: 5 * time.Second,
        },
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
