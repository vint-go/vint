---
title: useSliceAppend
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSliceAppend`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSliceAppend:
    # rule options here
```

## Details

Use a single `append` to concatenate two slices.

A `for` loop that appends each element from one slice to another can be replaced with a single `append` call using the `...` operator.

The rule only triggers when the range expression is a slice or array type. Ranging over maps or channels is not flagged, since the `...` spread operator does not work on those types.

Source: https://staticcheck.dev/docs/checks/#S1011

## Examples

### Invalid

```golang
package main

func merge(a, b []int) []int {
    for _, v := range b {
        a = append(a, v)
    }
    return a
}
```

### Valid

```golang
package main

func merge(a, b []int) []int {
    a = append(a, b...)
    return a
}
```

```golang
package main

// Ranging over a map cannot use the spread operator.
func collectValues(dst []int, m map[string]int) []int {
    for _, v := range m {
        dst = append(dst, v)
    }
    return dst
}
```
