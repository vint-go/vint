---
title: noRangeValCopy
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noRangeValCopy`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noRangeValCopy:
    sizeThreshold: 128   # size in bytes at which a range value copy is flagged (default: 128)
    skipTestFuncs: true   # whether to skip analysis of test functions (default: true)
```

## Details

Detects loops that copy large objects on each iteration. When a range loop captures values (e.g., `for _, x := range xs`), each iteration copies the value. For large structs, this copying is expensive. The checker suggests using index access or taking the address instead to avoid expensive copying operations.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
type BigStruct struct {
    Data [256]byte
}

var items []BigStruct
for _, item := range items { // copies BigStruct on each iteration
    process(item)
}
```

### Valid

```golang
type BigStruct struct {
    Data [256]byte
}

var items []BigStruct
for i := range items {
    process(&items[i]) // no copy, use index access
}
```

```golang
for i, item := range items {
    _ = i
    process(&item) // or take address of loop variable
}
```
