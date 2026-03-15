---
title: noHttpPost
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpPost`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpPost:
    # no additional options
```

## Details

Disallow calling `net/http.Post` without a context. The function `http.Post` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `net/http.NewRequestWithContext` to create a request with context, then execute it with `(*net/http.Client).Do`.

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
    // http.Post does not accept a context
    body := strings.NewReader(`{"key": "value"}`)
    resp, err := http.Post("https://example.com/api", "application/json", body)
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
