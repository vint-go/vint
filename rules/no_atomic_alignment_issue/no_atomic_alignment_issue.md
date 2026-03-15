---
title: noAtomicAlignmentIssue
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noAtomicAlignmentIssue`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noAtomicAlignmentIssue:
    # rule options here
```

## Details

Checks for non-64-bit-aligned arguments to `sync/atomic` functions. On 32-bit platforms, 64-bit atomic operations require that the variable be 64-bit aligned in memory. If a struct field used with `atomic.AddInt64`, `atomic.LoadInt64`, or similar functions is not properly aligned, it can cause a runtime panic on 32-bit architectures.

This analyzer helps catch alignment issues that would only manifest on 32-bit platforms, making them difficult to detect during development on 64-bit systems.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/atomicalign

## Examples

### Invalid

```golang
import "sync/atomic"

type Counter struct {
    flag bool       // 1 byte
    count int64     // not 64-bit aligned on 32-bit platforms due to preceding bool
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.count, 1) // may panic on 32-bit architectures
}
```

### Valid

```golang
import "sync/atomic"

type Counter struct {
    count int64     // 64-bit aligned (first field)
    flag  bool
}

func (c *Counter) Increment() {
    atomic.AddInt64(&c.count, 1) // safe: count is properly aligned
}
```
