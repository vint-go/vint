---
title: noRangeExprCopy
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noRangeExprCopy`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noRangeExprCopy:
    sizeThreshold: 512   # size in bytes at which a range expression copy is flagged (default: 512)
    skipTestFuncs: true   # whether to skip analysis of test functions (default: true)
```

## Details

Detects expensive copies of `for` loop range expressions. When iterating over an array in a for-range loop, Go copies the entire array. This checker warns when the copied data size reaches the configured threshold, suggesting the use of a pointer to the array (with `&`) to avoid the unnecessary copy.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var data [1024]int
for _, v := range data { // copies entire 1024-element array
    process(v)
}
```

### Valid

```golang
var data [1024]int
for _, v := range &data { // uses pointer, no copy
    process(v)
}
```

```golang
// Slices don't have this issue
var data []int
for _, v := range data {
    process(v)
}
```
