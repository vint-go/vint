---
title: noUnsignedNegativeComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnsignedNegativeComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnsignedNegativeComparison:
    # rule options here
```

## Details

Comparing unsigned values against negative values is pointless.

An unsigned integer can never be less than zero. Comparisons like `uint(x) < 0` are always false, and `uint(x) >= 0` are always true. This usually indicates a logic error.

Source: https://staticcheck.dev/docs/checks/#SA4003

## Examples

### Invalid

```golang
package main

func check(x uint) bool {
    // Always false: unsigned value cannot be negative
    return x < 0
}
```

### Valid

```golang
package main

func check(x uint) bool {
    // Meaningful comparison
    return x == 0
}
```
