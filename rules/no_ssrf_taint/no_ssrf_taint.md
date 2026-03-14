---
title: noSsrfTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noSsrfTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noSsrfTaint:
    # rule options here
```

## Details

Detects Server-Side Request Forgery (SSRF) vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources (such as HTTP request parameters) to HTTP request sinks (such as `http.Get`, `http.Post`, or `http.Client.Do`). It performs interprocedural analysis to detect cases where user-controlled URLs or URL components flow through multiple functions before being used in outbound HTTP requests.

SSRF attacks allow an attacker to make the server issue requests to arbitrary internal or external services. This can be used to bypass firewalls, access internal services (such as cloud metadata endpoints at `169.254.169.254`), scan internal networks, or exfiltrate data through DNS or HTTP channels.

URL inputs should be validated against an allowlist of permitted hosts and schemes. Internal and loopback addresses should be explicitly blocked.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func proxyHandler(w http.ResponseWriter, r *http.Request) {
    targetURL := r.URL.Query().Get("url")
    // Taint flows from request parameter to outbound HTTP request
    resp := fetchURL(targetURL)
    io.Copy(w, resp.Body)
}

func fetchURL(url string) *http.Response {
    resp, _ := http.Get(url)
    return resp
}
```

### Valid

```golang
var allowedHosts = map[string]bool{
    "api.example.com":     true,
    "cdn.example.com":     true,
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
    targetURL := r.URL.Query().Get("url")

    parsed, err := url.Parse(targetURL)
    if err != nil {
        http.Error(w, "invalid URL", http.StatusBadRequest)
        return
    }

    // Validate scheme
    if parsed.Scheme != "https" {
        http.Error(w, "only HTTPS allowed", http.StatusBadRequest)
        return
    }

    // Validate host against allowlist
    if !allowedHosts[parsed.Host] {
        http.Error(w, "host not allowed", http.StatusForbidden)
        return
    }

    // Block internal/loopback addresses
    ips, err := net.LookupIP(parsed.Hostname())
    if err != nil {
        http.Error(w, "DNS error", http.StatusBadGateway)
        return
    }
    for _, ip := range ips {
        if ip.IsLoopback() || ip.IsPrivate() {
            http.Error(w, "internal addresses not allowed", http.StatusForbidden)
            return
        }
    }

    resp, err := http.Get(targetURL)
    if err != nil {
        http.Error(w, "fetch error", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()
    io.Copy(w, resp.Body)
}
```
