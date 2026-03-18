---
title: noRedundantSliceExpression
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantSliceExpression`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantSliceExpression:
    # no additional options
```

## Details

Detects slice expressions that can be simplified to the expression itself. A slice expression like `s[:]` or `s[0:len(s)]` is equivalent to just `s` and adds unnecessary complexity. The redundant slicing should be removed for cleaner code.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
x := s[:]
```

```golang
x := s[0:len(s)]
```

### Valid

```golang
x := s
```
