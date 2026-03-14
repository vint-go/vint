---
title: noUnsafePackage
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noUnsafePackage`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noUnsafePackage:
    # rule options here
```

## Details

Audits the use of the `unsafe` package in Go code.

This rule detects calls to unsafe package functions including `unsafe.Pointer`, `unsafe.String`, `unsafe.StringData`, `unsafe.Slice`, and `unsafe.SliceData`. The `unsafe` package bypasses Go's type safety system and memory safety guarantees, which can lead to memory corruption, undefined behavior, and security vulnerabilities.

While there are legitimate use cases for the `unsafe` package (such as interfacing with C code or performance-critical operations), its usage should be carefully audited and minimized. Code using `unsafe` is not protected by the Go compatibility guarantee and may break with future Go versions.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "unsafe"

// Using unsafe.Pointer to cast between types
func unsafeCast(i *int) *float64 {
    return (*float64)(unsafe.Pointer(i))
}
```

```golang
import "unsafe"

// Using unsafe.Slice to create a slice from a pointer
func makeSlice(ptr *byte, length int) []byte {
    return unsafe.Slice(ptr, length)
}
```

```golang
import "unsafe"

// Using unsafe.String to create a string from bytes
func makeString(ptr *byte, length int) string {
    return unsafe.String(ptr, length)
}
```

### Valid

```golang
// Using encoding/binary for type conversion
import "encoding/binary"

func convertToFloat(data []byte) float64 {
    bits := binary.LittleEndian.Uint64(data)
    return math.Float64frombits(bits)
}
```

```golang
// Using standard library functions instead of unsafe
func copyBytes(src []byte) []byte {
    dst := make([]byte, len(src))
    copy(dst, src)
    return dst
}
```
