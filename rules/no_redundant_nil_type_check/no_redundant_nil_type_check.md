---
title: noRedundantNilTypeCheck
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantNilTypeCheck`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantNilTypeCheck:
    # rule options here
```

## Details

Omit redundant nil check in type assertion.

A two-value type assertion already returns `false` for nil interfaces. Checking `x != nil` before a type assertion is redundant.

Source: https://staticcheck.dev/docs/checks/#S1020

## Examples

### Invalid

```golang
package main

func process(x interface{}) {
    if x != nil {
        if v, ok := x.(string); ok {
            _ = v
        }
    }
}
```

### Valid

```golang
package main

func process(x interface{}) {
    if v, ok := x.(string); ok {
        _ = v
    }
}
```
