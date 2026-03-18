---
title: noWeakSliceGuard
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noWeakSliceGuard`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noWeakSliceGuard:
    # no additional options
```

## Details

Detects conditions that are unsafe due to not being exhaustive. This checker identifies two unsafe patterns where nil checks on slices are insufficient before indexing:

1. `x != nil && usageOf(x[i])` -- A nil check using inequality followed by array indexing.
2. `x == nil || usageOf(x[i])` -- An equality check with logical OR preceding indexing.

A nil slice can still panic when indexed if it is non-nil but empty. The checker recommends using `len(xs) != 0` instead of nil comparisons to properly guard against out-of-bounds access.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if xs != nil && xs[0] != nil {
    // xs could be non-nil but empty
}
```

### Valid

```golang
if len(xs) != 0 && xs[0] != nil {
    // properly checks length before indexing
}
```
