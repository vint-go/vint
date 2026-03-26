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

Detects slice expressions that can be simplified to the expression itself. A slice expression like `s[:]` or `s[0:len(s)]` is equivalent to just `s` when `s` is a slice or string, and adds unnecessary complexity. The redundant slicing should be removed for cleaner code.

Note: This rule does **not** flag `array[:]` expressions, because slicing an array (`[N]T`) produces a slice (`[]T`), which is a necessary type conversion. Removing `[:]` from an array would cause a compilation error.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
s := []int{1, 2, 3}
x := s[:]
```

```golang
s := []int{1, 2, 3}
x := s[0:len(s)]
```

```golang
s := "hello"
x := s[:]
```

### Valid

```golang
x := s
```

```golang
// array[:] is a necessary conversion from [N]T to []T
var a [5]int
x := a[:]
```
