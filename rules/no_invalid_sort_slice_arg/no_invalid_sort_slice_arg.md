---
title: noInvalidSortSliceArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidSortSliceArg`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidSortSliceArg:
    # rule options here
```

## Details

Checks for calls to `sort.Slice` that do not pass a slice type as the first argument. The `sort.Slice` function expects a slice as its first argument. Passing a non-slice type (such as a pointer to a slice or an array) will cause a runtime panic.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/sortslice

## Examples

### Invalid

```golang
import "sort"

func example() {
    s := []int{3, 1, 2}
    // Bad: passing a pointer to a slice instead of the slice itself
    sort.Slice(&s, func(i, j int) bool {
        return s[i] < s[j]
    })
}
```

### Valid

```golang
import "sort"

func example() {
    s := []int{3, 1, 2}
    // Good: passing the slice directly
    sort.Slice(s, func(i, j int) bool {
        return s[i] < s[j]
    })
}
```
