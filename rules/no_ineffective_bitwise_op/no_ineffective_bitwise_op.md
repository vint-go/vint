---
title: noIneffectiveBitwiseOp
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIneffectiveBitwiseOp`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIneffectiveBitwiseOp:
    # rule options here
```

## Details

Certain bitwise operations have no effect, such as `x ^ 0` or `x & 0`.

Operations like `x ^ 0` always equal `x`, `x & 0` always equals `0`, and `x | 0` always equals `x`. These are no-ops that indicate a logic error.

Source: https://staticcheck.dev/docs/checks/#SA4016

## Examples

### Invalid

```golang
package main

func process(x int) int {
    // x ^ 0 is always x
    return x ^ 0
}
```

### Valid

```golang
package main

func process(x int) int {
    // Meaningful bitwise operation
    return x ^ 0xFF
}
```
