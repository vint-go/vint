---
title: useCopyBuiltin
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCopyBuiltin`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCopyBuiltin:
    # rule options here
```

## Details

Replace for loop with call to `copy`.

A `for` loop that copies elements from one slice to another can be replaced with the built-in `copy` function, which is clearer and potentially more efficient.

Source: https://staticcheck.dev/docs/checks/#S1001

## Examples

### Invalid

```golang
package main

func clone(src []int) []int {
    dst := make([]int, len(src))
    for i, v := range src {
        dst[i] = v
    }
    return dst
}
```

### Valid

```golang
package main

func clone(src []int) []int {
    dst := make([]int, len(src))
    copy(dst, src)
    return dst
}
```
