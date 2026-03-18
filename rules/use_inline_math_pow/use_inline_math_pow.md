---
title: useInlineMathPow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useInlineMathPow`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useInlineMathPow:
    # rule options here
```

## Details

Expand call to `math.Pow`.

Calls to `math.Pow` with small integer exponents can be replaced with multiplication, which is more efficient and avoids floating-point conversion overhead.

Source: https://staticcheck.dev/docs/checks/#QF1005

## Examples

### Invalid

```golang
package main

import "math"

func square(x float64) float64 {
    return math.Pow(x, 2)
}
```

### Valid

```golang
package main

func square(x float64) float64 {
    return x * x
}
```
