---
title: noMissingReadHeaderTimeout
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noMissingReadHeaderTimeout`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noMissingReadHeaderTimeout:
    # rule options here
```

## Details

Detects when `ReadHeaderTimeout` is not configured on an `http.Server`, which can lead to Slowloris-type denial-of-service attacks.

The `ReadHeaderTimeout` field on `http.Server` specifies the maximum duration for reading request headers. If not set, the server will wait indefinitely for the client to send headers, allowing an attacker to hold connections open by slowly sending headers. This can exhaust the server's connection pool and prevent legitimate clients from connecting.

Slowloris attacks are particularly effective because they require minimal bandwidth from the attacker while consuming significant server resources. Setting `ReadHeaderTimeout` to a reasonable value (e.g., 5-10 seconds) mitigates this attack vector.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// Server without ReadHeaderTimeout
func main() {
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
    server.ListenAndServe()
}
```

```golang
import "net/http"

// Only setting ReadTimeout is not sufficient
func main() {
    server := &http.Server{
        Addr:        ":8080",
        Handler:     handler,
        ReadTimeout: 10 * time.Second,
    }
    server.ListenAndServe()
}
```

### Valid

```golang
import (
    "net/http"
    "time"
)

// Server with ReadHeaderTimeout configured
func main() {
    server := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadHeaderTimeout: 10 * time.Second,
    }
    server.ListenAndServe()
}
```

```golang
import (
    "net/http"
    "time"
)

// Server with both timeouts configured
func main() {
    server := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,
    }
    server.ListenAndServe()
}
```
