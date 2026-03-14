---
title: noServeWithoutTimeout
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noServeWithoutTimeout`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noServeWithoutTimeout:
    # rule options here
```

## Details

Detects the use of `net/http` serve functions that have no support for setting timeouts.

Functions like `http.ListenAndServe` and `http.ListenAndServeTLS` use a default `http.Server` with no timeout configuration. Without timeouts, the server is vulnerable to resource exhaustion attacks where clients can hold connections open indefinitely, consuming server resources.

Instead of using these convenience functions, create an explicit `http.Server` struct with appropriate timeout values (`ReadTimeout`, `ReadHeaderTimeout`, `WriteTimeout`, `IdleTimeout`) to protect against slow client attacks and resource exhaustion.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// Using ListenAndServe without timeout support
func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
```

```golang
import "net/http"

// Using ListenAndServeTLS without timeout support
func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
}
```

### Valid

```golang
import (
    "net/http"
    "time"
)

// Creating server with explicit timeouts
func main() {
    server := &http.Server{
        Addr:              ":8080",
        Handler:           handler,
        ReadTimeout:       10 * time.Second,
        ReadHeaderTimeout: 5 * time.Second,
        WriteTimeout:      10 * time.Second,
        IdleTimeout:       120 * time.Second,
    }
    server.ListenAndServe()
}
```
