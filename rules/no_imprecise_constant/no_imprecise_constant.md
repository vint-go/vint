---
title: noImpreciseConstant
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noImpreciseConstant`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noImpreciseConstant:
    # rule options here
```

## Details

Go constants always have an exact value, not an approximation.

Go constants are not IEEE 754 floating-point values. They have arbitrary precision. Writing something like `const x = 1.0 / 3.0` will give you the exact mathematical value, and it will only be truncated when assigned to a variable. There is no need to manually specify more decimal places.

Source: https://staticcheck.dev/docs/checks/#SA4026

## Examples

### Invalid

```golang
package main

// Unnecessary precision - Go handles this automatically
const pi = 3.14159265358979323846264338327950288419716939937510

func main() {
    var f float64 = pi
    _ = f
}
```

### Valid

```golang
package main

import "math"

func main() {
    // Use the standard library constant
    f := math.Pi
    _ = f
}
```
