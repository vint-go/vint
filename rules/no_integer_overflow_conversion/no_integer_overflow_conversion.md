---
title: noIntegerOverflowConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noIntegerOverflowConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noIntegerOverflowConversion:
    # rule options here
```

## Details

Detects type conversions that can lead to integer overflow.

This rule uses SSA (Static Single Assignment) analysis to identify integer type conversions where the source type has a larger range than the destination type. Converting a larger integer type to a smaller one (e.g., `int64` to `int32`, `uint64` to `uint32`) can cause silent data loss when the value exceeds the target type's range.

Integer overflow can lead to incorrect calculations, buffer overflows, denial of service, and security vulnerabilities. This is particularly dangerous when the converted value is used for memory allocation sizes, array indices, or security-critical calculations.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Converting int64 to int32 without bounds check
var bigVal int64 = math.MaxInt64
smallVal := int32(bigVal) // Overflow
```

```golang
// Converting uint to int may overflow on large values
func processSize(size uint) {
    intSize := int(size) // May overflow if size > MaxInt
    buffer := make([]byte, intSize)
    // ...
}
```

```golang
// Converting signed to unsigned without checking for negative
func toUint(val int) uint {
    return uint(val) // Negative values wrap around
}
```

### Valid

```golang
// Bounds checking before conversion
var bigVal int64 = getValue()
if bigVal > math.MaxInt32 || bigVal < math.MinInt32 {
    return fmt.Errorf("value %d out of int32 range", bigVal)
}
smallVal := int32(bigVal)
```

```golang
// Using math/bits for safe conversion
func safeUintToInt(val uint) (int, error) {
    if val > uint(math.MaxInt) {
        return 0, fmt.Errorf("value %d exceeds max int", val)
    }
    return int(val), nil
}
```

```golang
// Using explicit bit-sized types from the start
func processSize(size int32) {
    buffer := make([]byte, size)
    // ...
}
```
