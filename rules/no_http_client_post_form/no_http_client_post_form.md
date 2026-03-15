---
title: noHttpClientPostForm
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpClientPostForm`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpClientPostForm:
    # no additional options
```

## Details

Disallow calling `(*net/http.Client).PostForm` without a context. The method `(*http.Client).PostForm` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `(*net/http.Client).Do` with a request created via `net/http.NewRequestWithContext` instead.

Passing `context.Context` to HTTP operations enables the caller to cancel requests, propagate deadlines, and integrate with distributed tracing systems.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "net/http"
    "net/url"
)

func main() {
    client := &http.Client{}
    // (*http.Client).PostForm does not accept a context
    resp, err := client.PostForm("https://example.com/form", url.Values{
        "username": {"user"},
        "password": {"pass"},
    })
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
    "net/url"
    "strings"
)

func main() {
    ctx := context.Background()
    data := url.Values{
        "username": {"user"},
        "password": {"pass"},
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://example.com/form", strings.NewReader(data.Encode()))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println(resp.Status)
}
```
