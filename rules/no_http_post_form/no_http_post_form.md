---
title: noHttpPostForm
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttpPostForm`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttpPostForm:
    # no additional options
```

## Details

Disallow calling `net/http.PostForm` without a context. The function `http.PostForm` does not accept a `context.Context` parameter, which prevents the caller from controlling cancellation, deadlines, and tracing. Use `net/http.NewRequestWithContext` to create a request with context, then execute it with `(*net/http.Client).Do`.

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
    // http.PostForm does not accept a context
    resp, err := http.PostForm("https://example.com/form", url.Values{
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
