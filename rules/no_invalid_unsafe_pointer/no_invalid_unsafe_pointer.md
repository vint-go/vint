---
title: noInvalidUnsafePointer
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidUnsafePointer`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidUnsafePointer:
    # rule options here
```

## Details

Checks for invalid conversions of `uintptr` to `unsafe.Pointer`. The Go specification defines only a few legal patterns for converting between `uintptr` and `unsafe.Pointer`. Any conversion that does not follow these patterns is invalid and may break in future Go versions or on different platforms.

The legal patterns are:
1. Converting `unsafe.Pointer` to `uintptr` (for printing or arithmetic)
2. Converting `uintptr` back to `unsafe.Pointer` in the same expression
3. Using `unsafe.Pointer` with `reflect.Value.Pointer()` or `reflect.Value.UnsafeAddr()`
4. Passing `unsafe.Pointer` to `syscall.Syscall`

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/unsafeptr

## Examples

### Invalid

```golang
import "unsafe"

func example() {
    var x int
    // Bad: storing uintptr in a variable and converting later
    // The GC may move the object between the two lines
    ptr := uintptr(unsafe.Pointer(&x))
    p := unsafe.Pointer(ptr) // invalid: uintptr stored in variable
    _ = p
}
```

### Valid

```golang
import "unsafe"

func example() {
    var x int
    // Good: conversion in a single expression
    p := unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x))
    _ = p
}
```

```golang
import "unsafe"

func example() {
    var x int
    // Good: converting unsafe.Pointer to uintptr for printing
    fmt.Printf("address: %x\n", uintptr(unsafe.Pointer(&x)))
}
```
