---
title: useJoinHostPort
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useJoinHostPort`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useJoinHostPort:
    # rule options here
```

## Details

Checks for `net.Dial` calls that use IPv4-only address construction patterns such as `fmt.Sprintf("%s:%d", host, port)` or `fmt.Sprintf("%s:%s", host, port)`. These patterns do not work correctly with IPv6 addresses because IPv6 addresses contain colons and must be enclosed in square brackets in a host:port string (e.g., `[::1]:8080`).

Instead, use `net.JoinHostPort` which correctly handles both IPv4 and IPv6 addresses.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/hostport

## Examples

### Invalid

```golang
import (
    "fmt"
    "net"
)

func connect(host string, port int) (net.Conn, error) {
    // Bad: does not work with IPv6 addresses
    addr := fmt.Sprintf("%s:%d", host, port)
    return net.Dial("tcp", addr)
}
```

### Valid

```golang
import "net"

func connect(host string, port string) (net.Conn, error) {
    // Good: net.JoinHostPort handles IPv6 correctly
    addr := net.JoinHostPort(host, port)
    return net.Dial("tcp", addr)
}
```
