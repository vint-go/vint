---
title: noSuspiciousSortSlice
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSuspiciousSortSlice`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSuspiciousSortSlice:
    # no additional options
```

## Details

Detects suspicious `sort.Slice` calls. This checker identifies two problematic patterns in `sort.Slice` and `sort.SliceStable` calls:

1. **Missing slice reference** -- When the comparison function doesn't actually reference the slice being sorted, instead comparing values from a different collection.
2. **Reversed index parameters** -- When the comparison logic uses index `j` on the left operand and index `i` on the right operand, which may indicate a sorting logic error.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
sort.Slice(xs, func(i, j int) bool {
    return ys[i] < ys[j] // comparing ys instead of xs
})
```

```golang
sort.Slice(xs, func(i, j int) bool {
    return xs[j] < xs[i] // reversed indices
})
```

### Valid

```golang
sort.Slice(xs, func(i, j int) bool {
    return xs[i] < xs[j]
})
```
