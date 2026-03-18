---
title: useSimplifiedBoolReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSimplifiedBoolReturn`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSimplifiedBoolReturn:
    # rule options here
```

## Details

Simplify returning boolean expression.

An `if/else` that returns `true`/`false` based on a condition can be simplified to just returning the condition directly.

Source: https://staticcheck.dev/docs/checks/#S1008

## Examples

### Invalid

```golang
package main

func isPositive(x int) bool {
    if x > 0 {
        return true
    }
    return false
}
```

### Valid

```golang
package main

func isPositive(x int) bool {
    return x > 0
}
```
