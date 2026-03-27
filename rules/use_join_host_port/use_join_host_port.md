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

Checks for `fmt.Sprintf` calls that construct host:port addresses that do not work correctly with IPv6. It detects two scopes:

1. **URL-prefix patterns** like `fmt.Sprintf("http://%s:%d", host, port)` — these produce malformed URLs with IPv6 addresses.
2. **Bare host:port patterns** like `fmt.Sprintf("%s:%d", host, port)` — these produce ambiguous addresses with IPv6 hosts.

Use `net.JoinHostPort` instead, which correctly handles both IPv4 and IPv6 addresses.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/hostport

https://github.com/stbenjam/no-sprintf-host-port

## Examples

### Invalid

```golang
import "fmt"

func buildURL(host string, port int) string {
    // Bad: URL construction does not work with IPv6 addresses
    return fmt.Sprintf("http://%s:%d/path", host, port)
}

func dial(host string, port int) {
    // Bad: bare host:port does not work with IPv6 addresses
    addr := fmt.Sprintf("%s:%d", host, port)
    net.Dial("tcp", addr)
}
```

### Valid

```golang
import (
    "fmt"
    "net"
)

func buildURL(host string, port int) string {
    // Good: net.JoinHostPort handles IPv6 correctly
    return fmt.Sprintf("http://%s/path", net.JoinHostPort(host, fmt.Sprintf("%d", port)))
}

func listen(host string, port int) string {
    // Good: net.JoinHostPort handles IPv6 correctly
    addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
    return addr
}
```
