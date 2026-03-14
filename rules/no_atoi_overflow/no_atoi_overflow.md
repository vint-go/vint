---
title: noAtoiOverflow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noAtoiOverflow`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noAtoiOverflow:
    # rule options here
```

## Details

Detects potential integer overflow when converting `strconv.Atoi` results to smaller integer types such as `int32` or `int16`.

The `strconv.Atoi` function returns an `int`, which is 64 bits wide on most modern platforms. Converting this result directly to `int32` or `int16` without bounds checking can cause integer overflow, leading to unexpected values. This can result in security vulnerabilities when the converted value is used for buffer sizes, array indices, or other security-sensitive calculations.

Before casting the result of `strconv.Atoi` to a smaller integer type, the value should be validated to ensure it fits within the target type's range.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
// Direct conversion of Atoi result to int32 without bounds check
val, err := strconv.Atoi(input)
if err != nil {
    log.Fatal(err)
}
smallVal := int32(val) // Potential overflow
```

```golang
// Direct conversion to int16
val, _ := strconv.Atoi(input)
result := int16(val) // Potential overflow
```

### Valid

```golang
// Bounds checking before conversion
val, err := strconv.Atoi(input)
if err != nil {
    log.Fatal(err)
}
if val > math.MaxInt32 || val < math.MinInt32 {
    log.Fatal("value out of range for int32")
}
smallVal := int32(val)
```

```golang
// Using strconv.ParseInt with explicit bit size
val, err := strconv.ParseInt(input, 10, 32)
if err != nil {
    log.Fatal(err)
}
smallVal := int32(val)
```
