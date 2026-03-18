---
title: noMissingReturnAfterHttpError
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMissingReturnAfterHttpError`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMissingReturnAfterHttpError:
    # no additional options
```

## Details

Detects `http.Error` calls without a following `return` statement. When `http.Error` is called in an HTTP handler, a `return` should typically follow to prevent the handler from continuing to execute and potentially writing additional response data.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func handler(w http.ResponseWriter, r *http.Request) {
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        // missing return: handler continues executing
    }
    // ... more code that writes to w
}
```

### Valid

```golang
func handler(w http.ResponseWriter, r *http.Request) {
    if err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    // ... more code that writes to w
}
```
