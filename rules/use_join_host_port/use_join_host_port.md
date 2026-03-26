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

Checks for `fmt.Sprintf` calls that construct URLs with a scheme prefix and a host:port component, such as `fmt.Sprintf("http://%s:%d", host, port)` or `fmt.Sprintf("https://%s:%s", host, port)`. These patterns do not work correctly with IPv6 addresses because IPv6 addresses contain colons and must be enclosed in square brackets in a host:port string (e.g., `[::1]:8080`).

Bare host:port patterns like `fmt.Sprintf("%s:%d", host, port)` are **not** flagged, as they are commonly used for `net.Listen`, `http.Server.Addr`, etc., where the risk of IPv6 breakage is lower.

Instead, use `net.JoinHostPort` which correctly handles both IPv4 and IPv6 addresses.

Source: https://github.com/stbenjam/no-sprintf-host-port

## Examples

### Invalid

```golang
import "fmt"

func buildURL(host string, port int) string {
    // Bad: URL construction does not work with IPv6 addresses
    return fmt.Sprintf("http://%s:%d/path", host, port)
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
    // OK: bare host:port is acceptable for non-URL uses
    return fmt.Sprintf("%s:%d", host, port)
}
```
