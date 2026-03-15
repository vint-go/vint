---
title: noAtomicAssignMisuse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noAtomicAssignMisuse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noAtomicAssignMisuse:
    # rule options here
```

## Details

Checks for common mistakes using the `sync/atomic` package. Specifically, it detects cases where the result of an atomic operation (such as `atomic.AddInt64`) is not assigned back to the same variable that was passed as the argument. This often indicates a race condition where the programmer intended to atomically update a variable but instead discarded the result.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/atomic

## Examples

### Invalid

```golang
import "sync/atomic"

func example() {
    var x int64
    // Wrong: result of AddInt64 is not assigned back to x
    // This introduces a data race
    x = atomic.AddInt64(&x, 1)
}
```

### Valid

```golang
import "sync/atomic"

func example() {
    var x int64
    // Correct: use the return value directly
    newVal := atomic.AddInt64(&x, 1)
    _ = newVal
}
```

```golang
import "sync/atomic"

func example() {
    var x int64
    // Correct: simply call the function for its side effect
    atomic.AddInt64(&x, 1)
}
```
