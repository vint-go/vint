---
title: noExcessiveShift
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noExcessiveShift`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noExcessiveShift:
    # rule options here
```

## Details

Checks for shifts that exceed the width of an integer. Shifting an integer by more than or equal to its bit width produces a result that is always zero (for left shifts) or always zero/negative-one (for right shifts), which is almost certainly a programming error.

For example, shifting a `uint8` left by 8 or more bits always yields zero.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/shift

## Examples

### Invalid

```golang
func example() {
    var x uint8
    // Bad: shifting uint8 by 8 bits always yields 0
    y := x << 8
    _ = y
}
```

```golang
func example() {
    var x int32
    // Bad: shifting int32 by 32 bits exceeds its width
    y := x << 32
    _ = y
}
```

### Valid

```golang
func example() {
    var x uint8
    // Good: shift within the type's bit width
    y := x << 4
    _ = y
}
```

```golang
func example() {
    var x int64
    // Good: shift within the type's bit width
    y := x << 32
    _ = y
}
```
