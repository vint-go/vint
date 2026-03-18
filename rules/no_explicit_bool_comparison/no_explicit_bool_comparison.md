---
title: noExplicitBoolComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noExplicitBoolComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noExplicitBoolComparison:
    # rule options here
```

## Details

Omit comparison with boolean constant.

Comparing a boolean value to `true` or `false` is redundant. Use the boolean value directly or negate it.

Source: https://staticcheck.dev/docs/checks/#S1002

## Examples

### Invalid

```golang
package main

func check(b bool) {
    if b == true {
        // ...
    }
    if b == false {
        // ...
    }
}
```

### Valid

```golang
package main

func check(b bool) {
    if b {
        // ...
    }
    if !b {
        // ...
    }
}
```
