---
title: noUnsafeRedirectPolicy
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnsafeRedirectPolicy`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnsafeRedirectPolicy:
    # rule options here
```

## Details

Detects unsafe redirect policies that may propagate sensitive headers to different domains.

This rule uses SSA analysis to identify HTTP client configurations where the redirect policy (`CheckRedirect` function) does not strip sensitive headers (such as `Authorization`, cookies, or custom authentication headers) when following redirects to a different host. When an HTTP client follows a redirect from one domain to another, sensitive headers from the original request may be forwarded to the target domain, potentially leaking credentials or session tokens.

HTTP clients should implement redirect policies that remove sensitive headers when the redirect target is a different domain than the original request.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// Client that blindly follows redirects with all headers
client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        // Copying all headers including Authorization to redirect target
        for key, val := range via[0].Header {
            req.Header[key] = val
        }
        return nil
    },
}

req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
req.Header.Set("Authorization", "Bearer secret-token")
client.Do(req) // Token may leak if redirected to another domain
```

### Valid

```golang
import "net/http"

// Client that strips sensitive headers on cross-domain redirects
client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        if len(via) > 0 && req.URL.Host != via[0].URL.Host {
            // Strip sensitive headers when redirecting to different host
            req.Header.Del("Authorization")
            req.Header.Del("Cookie")
        }
        return nil
    },
}
```

```golang
import "net/http"

// Disabling automatic redirects entirely
client := &http.Client{
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },
}
```
