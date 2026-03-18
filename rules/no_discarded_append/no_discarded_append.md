---
title: noDiscardedAppend
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDiscardedAppend`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDiscardedAppend:
    # rule options here
```

## Details

The result of `append` is never used, except maybe in other calls to `append`.

The `append` function returns a new slice. If the return value is not assigned to anything, the appended elements are lost. This is almost always a bug.

Source: https://staticcheck.dev/docs/checks/#SA4010

## Examples

### Invalid

```golang
package main

func process() []int {
    s := []int{1, 2, 3}
    // Return value of append is discarded
    append(s, 4)
    return s
}
```

### Valid

```golang
package main

func process() []int {
    s := []int{1, 2, 3}
    // Assign the result back
    s = append(s, 4)
    return s
}
```
