---
title: noHttpGet
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpGet`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpGet:
    # no additional options
```

## Details

Disallow calling `net/http.Get` without a context. The function `http.Get` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `net/http.NewRequestWithContext` to create a request with context, then execute it with `(*net/http.Client).Do`.

Passing `context.Context` to HTTP operations enables the caller to cancel requests, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // http.Get does not accept a context
    resp, err := http.Get("https://example.com")
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
)

func main() {
    ctx := context.Background()
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
    if err != nil {
        panic(err)
    }
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println(resp.Status)
}
```
