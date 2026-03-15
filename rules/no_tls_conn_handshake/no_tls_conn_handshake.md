---
title: noTlsConnHandshake
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTlsConnHandshake`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTlsConnHandshake:
    # no additional options
```

## Details

Disallow calling `(*crypto/tls.Conn).Handshake` without a context. The method `(*tls.Conn).Handshake` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*crypto/tls.Conn).HandshakeContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` to TLS handshake operations enables the caller to cancel slow handshakes, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "crypto/tls"
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "example.com:443")
    if err != nil {
        panic(err)
    }
    tlsConn := tls.Client(conn, &tls.Config{})
    // (*tls.Conn).Handshake does not accept a context
    err = tlsConn.Handshake()
    if err != nil {
        panic(err)
    }
    defer tlsConn.Close()
    fmt.Println("Handshake completed")
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
)

func main() {
    ctx := context.Background()
    conn, err := net.Dial("tcp", "example.com:443")
    if err != nil {
        panic(err)
    }
    tlsConn := tls.Client(conn, &tls.Config{})
    err = tlsConn.HandshakeContext(ctx)
    if err != nil {
        panic(err)
    }
    defer tlsConn.Close()
    fmt.Println("Handshake completed")
}
```
