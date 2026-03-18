---
title: noNilContext
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNilContext`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNilContext:
    # rule options here
```

## Details

A nil `context.Context` is being passed to a function, use `context.TODO` instead.

Passing a `nil` context violates the `context.Context` contract. Functions that accept a context expect a non-nil value. If you don't have a context available, use `context.TODO()` to indicate that the context should be replaced in the future, or `context.Background()` for a top-level context.

Source: https://staticcheck.dev/docs/checks/#SA1012

## Examples

### Invalid

```golang
package main

import (
    "context"
    "net/http"
)

func main() {
    // Passing nil context
    req, _ := http.NewRequestWithContext(nil, "GET", "http://example.com", nil)
    _ = req
}
```

### Valid

```golang
package main

import (
    "context"
    "net/http"
)

func main() {
    // Using context.Background() or context.TODO()
    req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
    _ = req
}
```
