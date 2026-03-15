---
title: noHttpClientPost
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpClientPost`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpClientPost:
    # no additional options
```

## Details

Disallow calling `(*net/http.Client).Post` without a context. The method `(*http.Client).Post` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net/http.Client).Do` with a request created via `net/http.NewRequestWithContext` instead.

Passing `context.Context` to HTTP operations enables the caller to cancel requests, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "net/http"
    "strings"
)

func main() {
    client := &http.Client{}
    body := strings.NewReader(`{"key": "value"}`)
    // (*http.Client).Post does not accept a context
    resp, err := client.Post("https://example.com/api", "application/json", body)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println(resp.Status)
}
```

### Valid

```golang
package main

import (
    "context"
    "fmt"
    "net/http"
    "strings"
)

func main() {
    ctx := context.Background()
    body := strings.NewReader(`{"key": "value"}`)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/api", body)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println(resp.Status)
}
```
