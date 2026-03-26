---
title: noSsrfViaVariable
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSsrfViaVariable`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSsrfViaVariable:
    # rule options here
```

## Details

Detects potential Server-Side Request Forgery (SSRF) vulnerabilities by identifying HTTP requests made with URLs derived from variables.

This rule monitors calls to `net/http` package-level functions `Get`, `Head`, `Post`, and `PostForm`, matching the scope of gosec rule G107. It flags cases where the URL argument is a variable identifier (not resolvable to a constant). String literals, constant identifiers, call expressions (e.g., `fmt.Sprintf(...)`), and binary expressions are not flagged.

Note: `http.NewRequest` and `http.NewRequestWithContext` are not monitored, as they are not in the scope of gosec G107.

SSRF attacks can be used to bypass firewalls, access internal services, scan internal networks, and exfiltrate data. URL inputs should be validated against an allowlist of permitted hosts or URL patterns.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// URL from variable - potential SSRF
func fetch(url string) (*http.Response, error) {
    return http.Get(url)
}
```

```golang
import "net/http"

// URL constructed from user input, stored in variable
func fetchUser(userID string) (*http.Response, error) {
    url := "https://api.example.com/users/" + userID
    return http.Get(url)
}
```

### Valid

```golang
import "net/http"

// Hardcoded URL constant - safe
func fetchConfig() (*http.Response, error) {
    return http.Get("https://config.internal.example.com/settings")
}
```

```golang
import "net/http"

// Constant identifier - safe
const apiURL = "https://api.example.com/data"

func fetchData() (*http.Response, error) {
    return http.Get(apiURL)
}
```

```golang
import (
    "fmt"
    "net/http"
)

// Call expression (fmt.Sprintf) - not flagged
func fetchResource(host string) (*http.Response, error) {
    return http.Get(fmt.Sprintf("https://%s/api/resource", host))
}
```

```golang
import "net/http"

// URL validated against allowlist before use
func fetchValidated(rawURL string) (*http.Response, error) {
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return nil, err
    }
    allowedHosts := map[string]bool{"api.example.com": true}
    if !allowedHosts[parsed.Host] {
        return nil, fmt.Errorf("host not allowed: %s", parsed.Host)
    }
    return http.Get(rawURL)
}
```
