---
title: noHttpNewRequest
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpNewRequest`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpNewRequest:
    # no additional options
```

## Details

Disallow calling `net/http.NewRequest` without a context. The function `http.NewRequest` does not accept a `context.Context` parameter, which means the resulting request will use `context.Background()` by default. Use `net/http.NewRequestWithContext` instead, which accepts a `context.Context` as its first argument.

Passing `context.Context` when creating HTTP requests enables the caller to cancel requests, propagate deadlines, and integrate with distributed tracing systems.

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
    // http.NewRequest does not accept a context
    req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
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
