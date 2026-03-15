---
title: noUnsafeCorsBypass
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnsafeCorsBypass`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnsafeCorsBypass:
    # rule options here
```

## Details

Detects unsafe Cross-Origin Resource Protection (CORP) bypass patterns.

This rule uses SSA analysis to identify configurations or code patterns that improperly relax cross-origin protections, such as setting overly permissive CORS (Cross-Origin Resource Sharing) headers or disabling cross-origin checks entirely. Misconfigured CORS policies can allow malicious websites to make authenticated requests to your application on behalf of users, potentially leading to data theft or unauthorized actions.

Cross-origin policies should be configured with explicit, restrictive allowlists of permitted origins rather than using wildcard values or reflecting the request origin without validation.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    // Reflecting the request origin without validation
    origin := r.Header.Get("Origin")
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}
```

```golang
import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    // Using wildcard with credentials
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
```

### Valid

```golang
import "net/http"

var allowedOrigins = map[string]bool{
    "https://app.example.com":  true,
    "https://admin.example.com": true,
}

func handler(w http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
    if allowedOrigins[origin] {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Credentials", "true")
    }
}
```
