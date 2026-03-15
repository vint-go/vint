---
title: noMismatchedAppendAssign
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noMismatchedAppendAssign`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noMismatchedAppendAssign:
    # no additional options
```

## Details

Detects suspicious append result assignments. This checker identifies cases where an `append` operation assigns its result to a different slice than the one being appended to. This is usually a copy-paste error or a logic mistake.

The checker compares the assignment target (left side) with the first argument to `append()` (the slice being extended). It also handles exceptions such as blank identifier assignments, ellipsis patterns, and index expressions with map/slice operations.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// The result of append is assigned to a different slice
p.positives = append(p.negatives, x)
```

```golang
// Appending to 'b' but assigning to 'a'
a = append(b, item)
```

### Valid

```golang
// The result of append is assigned to the same slice
p.positives = append(p.positives, x)
```

```golang
// Using ellipsis pattern is allowed
xs = append(ys, xs...)
```
