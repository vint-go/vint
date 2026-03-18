---
title: useHttpNoBody
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useHttpNoBody`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useHttpNoBody:
    # no additional options
```

## Details

Detects nil usages in `http.NewRequest` calls, suggesting `http.NoBody` instead. When creating an HTTP request without a body, using `http.NoBody` is more explicit and idiomatic than passing `nil`. The `http.NoBody` constant was introduced in Go 1.8 specifically for this purpose.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
req, err := http.NewRequest("GET", url, nil)
```

### Valid

```golang
req, err := http.NewRequest("GET", url, http.NoBody)
```
