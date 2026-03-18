---
title: useDirectStringComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useDirectStringComparison`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useDirectStringComparison:
    # no additional options
```

## Details

Detects empty string checks using `len(s) == 0` that can be written more idiomatically as `s == ""`. Direct string comparison against the empty string literal is clearer in intent and is the preferred Go idiom for testing whether a string is empty.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if len(s) == 0 {
    handleEmpty()
}
```

### Valid

```golang
if s == "" {
    handleEmpty()
}
```

```golang
if s != "" {
    handleNonEmpty()
}
```
