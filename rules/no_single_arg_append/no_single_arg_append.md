---
title: noSingleArgAppend
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noSingleArgAppend`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noSingleArgAppend:
    # rule options here
```

## Details

Detects if there is only one variable in an `append` call. This analyzer checks for cases where `append` is called with only a single argument (the slice itself), which has no effect and is almost certainly a mistake.

A common error is writing `append(s)` instead of `append(s, elem)`, which does not modify the slice and indicates a missing element argument.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/appends

## Examples

### Invalid

```golang
// append with only one argument has no effect
func example() {
    s := []int{1, 2, 3}
    s = append(s) // missing element to append
}
```

### Valid

```golang
// append with proper arguments
func example() {
    s := []int{1, 2, 3}
    s = append(s, 4)
}
```

```golang
// append with variadic expansion
func example() {
    s := []int{1, 2, 3}
    t := []int{4, 5}
    s = append(s, t...)
}
```
