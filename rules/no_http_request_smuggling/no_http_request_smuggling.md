---
title: noHttpRequestSmuggling
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noHttpRequestSmuggling`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noHttpRequestSmuggling:
    # rule options here
```

## Details

Detects potential HTTP request smuggling vulnerabilities caused by conflicting headers or bare line feed (LF) characters in HTTP requests.

HTTP request smuggling occurs when front-end and back-end servers interpret the boundaries between HTTP requests differently. This can happen when conflicting `Content-Length` and `Transfer-Encoding` headers are present, or when bare LF characters (without carriage return) are used as line terminators. Different servers may parse these ambiguities differently, allowing attackers to "smuggle" requests through intermediate proxies.

This rule uses SSA (Static Single Assignment) analysis to detect patterns that could lead to request smuggling attacks. Proper HTTP implementations should reject requests with conflicting headers and require CRLF line endings.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "net/http"

// Constructing HTTP request with bare LF characters
func buildRequest(body string) string {
    return "GET / HTTP/1.1\n" +
        "Host: example.com\n" +
        "Content-Length: " + strconv.Itoa(len(body)) + "\n\n" +
        body
}
```

```golang
// Setting conflicting headers
req, _ := http.NewRequest("POST", url, body)
req.Header.Set("Content-Length", "10")
req.Header.Set("Transfer-Encoding", "chunked")
```

### Valid

```golang
import "net/http"

// Using standard library's HTTP client which handles headers correctly
func makeRequest(url string) (*http.Response, error) {
    return http.Get(url)
}
```

```golang
import "net/http"

// Using NewRequest with proper body handling
func makePost(url string, data io.Reader) (*http.Response, error) {
    req, err := http.NewRequest("POST", url, data)
    if err != nil {
        return nil, err
    }
    return http.DefaultClient.Do(req)
}
```
