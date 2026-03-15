---
title: noFramePointerClobber
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noFramePointerClobber`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noFramePointerClobber:
    # rule options here
```

## Details

Reports assembly code that clobbers the frame pointer before saving it. On architectures that use a frame pointer (such as amd64), assembly functions must save the frame pointer (BP register) before modifying it. Failing to do so breaks stack unwinding, which affects debugging, profiling, and panic stack traces.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/framepointer

## Examples

### Invalid

```asm
// Bad: clobbers BP without saving it first
TEXT ·example(SB), NOSPLIT, $0-8
    MOVQ $0, BP    // clobbers frame pointer
    RET
```

### Valid

```asm
// Good: saves and restores BP
TEXT ·example(SB), $8-8
    MOVQ BP, saved-8(SP)   // save frame pointer
    LEAQ saved-8(SP), BP   // set new frame pointer
    // ... function body ...
    MOVQ saved-8(SP), BP   // restore frame pointer
    RET
```
