---
title: noTruncatingComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noTruncatingComparison`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noTruncatingComparison:
    skipArchDependent: true  # skip architecture-dependent types like int, uint, uintptr (default: true)
```

## Details

Detects potential truncation issues when comparing ints of different sizes. This checker identifies problematic comparisons where an integer is cast to a smaller type before being compared with another integer. Casting the narrower operand to the larger type would be safer and more correct.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Casting x (int32) to smaller int16 before comparison
if int16(x) < y {
    // potential truncation
}
```

### Valid

```golang
// Cast the narrower operand to the larger type instead
if x < int32(y) {
    // safe comparison
}
```
