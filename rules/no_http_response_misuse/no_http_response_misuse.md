---
title: noHttpResponseMisuse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpResponseMisuse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpResponseMisuse:
    # rule options here
```

## Details

Checks for mistakes using HTTP responses. A common mistake is to defer the closing of the response body before checking the error returned by `http.Client.Do` or similar methods. If the error is non-nil, the response may be nil, causing a nil pointer dereference in the deferred `Body.Close()` call.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/httpresponse

## Examples

### Invalid

```golang
import "net/http"

func example() {
    resp, err := http.Get("https://example.com")
    // Bad: deferring Body.Close before checking error
    // If err != nil, resp may be nil
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
}
```

### Valid

```golang
import "net/http"

func example() {
    resp, err := http.Get("https://example.com")
    if err != nil {
        log.Fatal(err)
    }
    // Good: check error before deferring Body.Close
    defer resp.Body.Close()
}
```
