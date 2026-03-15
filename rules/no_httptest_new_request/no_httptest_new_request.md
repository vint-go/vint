---
title: noHttptestNewRequest
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noHttptestNewRequest`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noHttptestNewRequest:
    # no additional options
```

## Details

Disallow calling `net/http/httptest.NewRequest` without a context. The function `httptest.NewRequest` does not accept a `context.Context` parameter, which means the resulting request will use `context.Background()` by default. Use `net/http/httptest.NewRequestWithContext` instead, which accepts a `context.Context` as its first argument.

Even in tests, passing `context.Context` when creating HTTP requests enables proper cancellation, deadline propagation, and tracing integration, which can be important for testing context-aware handlers.

Source: https://github.com/sonatard/noctx

## Examples

### Invalid

```golang
package main_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    // httptest.NewRequest does not accept a context
    req := httptest.NewRequest(http.MethodGet, "/path", nil)
    w := httptest.NewRecorder()
    handler(w, req)
}
```

### Valid

```golang
package main_test

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    ctx := context.Background()
    req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/path", nil)
    w := httptest.NewRecorder()
    handler(w, req)
}
```
