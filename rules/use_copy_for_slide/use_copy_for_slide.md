---
title: useCopyForSlide
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCopyForSlide`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCopyForSlide:
    # rule options here
```

## Details

Use `copy` for sliding elements.

A loop that shifts elements within a slice can be replaced with a `copy` call, which is more efficient and idiomatic.

Source: https://staticcheck.dev/docs/checks/#S1018

## Examples

### Invalid

```golang
package main

func shift(s []int) {
    for i := 0; i < len(s)-1; i++ {
        s[i] = s[i+1]
    }
}
```

### Valid

```golang
package main

func shift(s []int) {
    copy(s, s[1:])
}
```
