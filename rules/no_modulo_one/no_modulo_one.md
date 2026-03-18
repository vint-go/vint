---
title: noModuloOne
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noModuloOne`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noModuloOne:
    # rule options here
```

## Details

`x % 1` is always zero.

The modulo operation `x % 1` always results in zero for any integer `x`. This is likely a logic error; perhaps you meant `x % 2` or another divisor.

Source: https://staticcheck.dev/docs/checks/#SA4028

## Examples

### Invalid

```golang
package main

func isOdd(x int) bool {
    // x % 1 is always 0
    return x%1 != 0
}
```

### Valid

```golang
package main

func isOdd(x int) bool {
    return x%2 != 0
}
```
