---
title: noRedundantVarType
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantVarType`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantVarType:
    # rule options here
```

## Details

Redundant type in variable declaration.

When the type of a variable can be inferred from the right-hand side, specifying it explicitly is redundant.

Source: https://staticcheck.dev/docs/checks/#ST1023

## Examples

### Invalid

```golang
package main

func process() {
    var x int = 42
    _ = x
}
```

### Valid

```golang
package main

func process() {
    var x = 42
    _ = x
    // Or better:
    y := 42
    _ = y
}
```
