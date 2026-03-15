---
title: noAsmDeclMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noAsmDeclMismatch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noAsmDeclMismatch:
    # rule options here
```

## Details

Reports mismatches between assembly (`.s`) files and Go declarations. This analyzer checks that assembly functions match their Go function prototypes, ensuring that parameter sizes, offsets, and return values are consistent between the Go declaration and the assembly implementation.

Common issues detected include:
- Mismatched argument sizes between Go and assembly
- Wrong frame size in assembly TEXT directives
- Incorrect offsets for function parameters and return values

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/asmdecl

## Examples

### Invalid

```golang
// Go declaration:
package example

func Add(x, y int64) int64
```

```asm
// Assembly file with wrong frame size:
TEXT ·Add(SB), NOSPLIT, $0-16  // wrong: should be $0-24 for two int64 args + int64 return
    MOVQ x+0(FP), AX
    ADDQ y+8(FP), AX
    MOVQ AX, ret+16(FP)
    RET
```

### Valid

```golang
// Go declaration:
package example

func Add(x, y int64) int64
```

```asm
// Assembly with correct frame size:
TEXT ·Add(SB), NOSPLIT, $0-24  // correct: 8+8+8 = 24 bytes
    MOVQ x+0(FP), AX
    ADDQ y+8(FP), AX
    MOVQ AX, ret+16(FP)
    RET
```
