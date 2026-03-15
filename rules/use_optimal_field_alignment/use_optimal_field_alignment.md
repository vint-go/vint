---
title: useOptimalFieldAlignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useOptimalFieldAlignment`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useOptimalFieldAlignment:
    # rule options here
```

## Details

Detects structs that would use less memory if their fields were sorted. Due to memory alignment requirements, the order of fields in a struct can affect its total size. The Go compiler inserts padding between fields to satisfy alignment constraints. By reordering fields from largest to smallest, the amount of padding can be minimized, reducing the struct's memory footprint.

This analyzer reports structs where reordering the fields would save memory and suggests the optimal field order.

Note: This analyzer is not enabled by default because field order sometimes carries semantic meaning or improves readability.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/fieldalignment

## Examples

### Invalid

```golang
// struct is 24 bytes due to padding
type Inefficient struct {
    a bool    // 1 byte + 7 bytes padding
    b int64   // 8 bytes
    c bool    // 1 byte + 7 bytes padding
}
```

### Valid

```golang
// struct is 16 bytes with optimized field order
type Efficient struct {
    b int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte + 6 bytes padding
}
```
