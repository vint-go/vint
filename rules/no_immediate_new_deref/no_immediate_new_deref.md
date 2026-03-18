---
title: noImmediateNewDeref
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noImmediateNewDeref`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noImmediateNewDeref:
    # no additional options
```

## Details

Detects immediate dereferencing of `new` expressions. This checker identifies code patterns where a `new()` call is immediately dereferenced with the `*` operator. Using the zero value directly is cleaner and more idiomatic. The checker makes an exception for type parameters, allowing `*new(T)` when T is a generic type parameter.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
x := *new(bool)
```

```golang
y := *new(int)
```

### Valid

```golang
x := false
```

```golang
y := 0
```

```golang
// Allowed for generic type parameters
func zero[T any]() T {
    return *new(T)
}
```
