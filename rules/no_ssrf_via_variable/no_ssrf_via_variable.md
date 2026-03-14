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

Detects potential Server-Side Request Forgery (SSRF) vulnerabilities by identifying HTTP requests made with URLs derived from variables or user input.

This rule monitors calls to `net/http` methods including `Do`, `Get`, `Head`, `Post`, `PostForm`, and `RoundTrip`, and flags cases where the URL argument is sourced from a variable rather than a hardcoded constant string. Variable-sourced URLs could be attacker-controlled, enabling SSRF attacks that allow an attacker to make the server issue requests to arbitrary internal or external services.

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

// URL constructed from user input
func fetchUser(userID string) (*http.Response, error) {
    url := "https://api.example.com/users/" + userID
    return http.Get(url)
}
```

```golang
import "net/http"

// URL from fmt.Sprintf
func fetchResource(host string) (*http.Response, error) {
    url := fmt.Sprintf("https://%s/api/resource", host)
    req, _ := http.NewRequest("GET", url, nil)
    return http.DefaultClient.Do(req)
}
```

### Valid

```golang
import "net/http"

// Hardcoded URL constant
func fetchConfig() (*http.Response, error) {
    return http.Get("https://config.internal.example.com/settings")
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
