---
title: noNanComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNanComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNanComparison:
    # rule options here
```

## Details

Comparing a value against NaN even though no value is equal to NaN.

In IEEE 754 floating-point arithmetic, NaN is not equal to anything, including itself. Comparisons like `x == math.NaN()` will always be false. Use `math.IsNaN(x)` instead.

Source: https://staticcheck.dev/docs/checks/#SA4012

## Examples

### Invalid

```golang
package main

import "math"

func isNaN(x float64) bool {
    // Always false: NaN != NaN
    return x == math.NaN()
}
```

### Valid

```golang
package main

import "math"

func isNaN(x float64) bool {
    return math.IsNaN(x)
}
```
