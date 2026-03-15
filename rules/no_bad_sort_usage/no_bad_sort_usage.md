---
title: noBadSortUsage
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noBadSortUsage`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noBadSortUsage:
    # no additional options
```

## Details

Detects bad usage of the `sort` package. This checker identifies patterns where sort functions are called incorrectly, such as sorting a slice with the wrong comparison function or using sort methods that don't match the intended ordering.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
sort.Ints(xs)
// xs is modified, but then used unsorted elsewhere
```

### Valid

```golang
sort.Ints(xs)
// use xs as sorted
```
