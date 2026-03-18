---
title: noImpossibleBuiltinResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noImpossibleBuiltinResult`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noImpossibleBuiltinResult:
    # rule options here
```

## Details

Checking for impossible return value from a builtin function.

Some builtins have guaranteed return value constraints. For example, `len()` and `cap()` always return non-negative values. Checking `len(x) < 0` is always false.

Source: https://staticcheck.dev/docs/checks/#SA4024

## Examples

### Invalid

```golang
package main

func process(s []int) {
    // len() never returns negative values
    if len(s) < 0 {
        panic("impossible")
    }
}
```

### Valid

```golang
package main

func process(s []int) {
    if len(s) == 0 {
        return
    }
}
```
