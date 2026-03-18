---
title: noUselessMathCall
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUselessMathCall`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUselessMathCall:
    # rule options here
```

## Details

Calling functions like `math.Ceil` on floats converted from integers doesn't do anything useful.

Functions like `math.Ceil`, `math.Floor`, and `math.Round` have no effect on values that are already integers. Converting an `int` to `float64` and then calling these functions is pointless.

Source: https://staticcheck.dev/docs/checks/#SA4015

## Examples

### Invalid

```golang
package main

import "math"

func process(x int) float64 {
    // Ceil of an integer value is the same value
    return math.Ceil(float64(x))
}
```

### Valid

```golang
package main

func process(x int) float64 {
    return float64(x)
}
```
