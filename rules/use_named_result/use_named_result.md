---
title: useNamedResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useNamedResult`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useNamedResult:
    checkExported: false  # whether to check exported functions (default: false)
```

## Details

Detects unnamed results that may benefit from names. This checker identifies function return values lacking names that could benefit from being named to improve code clarity.

For 2-return functions, it warns unless both returns share the same type (unless the second is `error` or `bool`). For functions returning more than 2 values, it warns when duplicate types appear, except for trailing `error` or `bool` returns.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f() (float64, float64) {
    return 0.0, 0.0
}
```

### Valid

```golang
func f() (x, y float64) {
    return 0.0, 0.0
}
```

```golang
func g() (int, error) {
    return 0, nil // trailing error is excluded
}
```
