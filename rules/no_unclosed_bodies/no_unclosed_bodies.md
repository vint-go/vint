---
title: noUnclosedBodies
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnclosedBodies`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnclosedBodies:
    # rule options here
```

## Details

Checks whether HTTP response bodies are correctly closed after use.

When making HTTP requests in Go, the resulting `http.Response.Body` must be explicitly closed. Failing to close the response body prevents the underlying TCP connection from being reused for subsequent requests, which can lead to resource leaks, exhausted file descriptors, and degraded application performance.

The standard pattern is to defer closing the body immediately after checking for errors:

```golang
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
```

This rule detects cases where `resp.Body.Close()` is never called on an HTTP response, including when using `http.Get`, `http.Post`, `http.Do`, and other HTTP client methods.

Source: https://github.com/timakin/bodyclose

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "io"
    "net/http"
)

func fetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    // resp.Body is never closed!
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return body, nil
}
```

```golang
package main

import "net/http"

func makeRequest(client *http.Client, req *http.Request) error {
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    // resp.Body is never closed!
    _ = resp.StatusCode
    return nil
}
```

```golang
package main

import "net/http"

func discardResponse(url string) error {
    // Response assigned to blank identifier — body can never be closed!
    _, err := http.Get(url)
    return err
}
```

### Valid

```golang
package main

import (
    "io"
    "net/http"
)

func fetchData(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    return body, nil
}
```

```golang
package main

import "net/http"

func makeRequest(client *http.Client, req *http.Request) error {
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    _ = resp.StatusCode
    return nil
}
```
