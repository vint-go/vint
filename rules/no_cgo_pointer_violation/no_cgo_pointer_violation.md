---
title: noCgoPointerViolation
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noCgoPointerViolation`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noCgoPointerViolation:
    # rule options here
```

## Details

Detects violations of the cgo pointer passing rules. When calling C functions from Go, there are strict rules about which Go pointers may be passed to C code. In general, Go code may pass a Go pointer to C provided the Go memory to which it points does not contain any Go pointers.

This analyzer checks that cgo calls do not pass Go pointers that contain other Go pointers, which would violate the cgo safety rules and could lead to runtime panics or memory corruption.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/cgocall

## Examples

### Invalid

```golang
/*
#include <stdlib.h>

void process(void *p) {}
*/
import "C"
import "unsafe"

func example() {
    s := "hello"
    // Passing a pointer to a Go string (which contains a Go pointer)
    C.process(unsafe.Pointer(&s)) // violates cgo pointer passing rules
}
```

### Valid

```golang
/*
#include <stdlib.h>

void process(void *p) {}
*/
import "C"
import "unsafe"

func example() {
    b := make([]byte, 5)
    copy(b, "hello")
    // Passing a pointer to a byte slice's underlying array (no nested Go pointers)
    C.process(unsafe.Pointer(&b[0]))
}
```
