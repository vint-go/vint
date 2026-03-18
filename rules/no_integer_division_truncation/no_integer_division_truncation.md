---
title: noIntegerDivisionTruncation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIntegerDivisionTruncation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIntegerDivisionTruncation:
    # rule options here
```

## Details

Integer division of literals that results in zero.

Dividing two integer literals where the numerator is smaller than the denominator always results in zero. This usually indicates that floating-point division was intended.

Source: https://staticcheck.dev/docs/checks/#SA4025

## Examples

### Invalid

```golang
package main

func main() {
    // Integer division: 1/2 = 0
    ratio := 1 / 2
    _ = ratio
}
```

### Valid

```golang
package main

func main() {
    // Floating-point division
    ratio := 1.0 / 2.0
    _ = ratio
}
```
