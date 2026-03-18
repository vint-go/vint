---
title: noUnnecessaryDeref
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryDeref`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryDeref:
    skipRecvDeref: true  # skip dereferences on pointer receiver method calls (default: true)
```

## Details

Detects dereference expressions that can be omitted. Go automatically dereferences pointers for field access and array indexing, so explicit dereferencing is unnecessary in many cases. This checker identifies two patterns:

1. **Selector expressions:** Cases like `(*k).field` that should be written as `k.field`.
2. **Index expressions:** Cases like `(*a)[5]` (where `a` is a pointer to an array) that should be written as `a[5]`.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
(*k).field = value
```

```golang
v := (*arr)[5]
```

### Valid

```golang
k.field = value
```

```golang
v := arr[5]
```
