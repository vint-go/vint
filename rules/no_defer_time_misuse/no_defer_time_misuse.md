---
title: noDeferTimeMisuse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeferTimeMisuse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeferTimeMisuse:
    # rule options here
```

## Details

Checks for common mistakes in defer statements. This analyzer detects problematic patterns in deferred function calls, such as deferring a call to `time.Since` which evaluates its argument immediately rather than at defer time, making it useless for measuring elapsed time.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/defers

## Examples

### Invalid

```golang
import (
    "log"
    "time"
)

func example() {
    start := time.Now()
    // Bad: time.Since(start) is evaluated immediately, not when deferred
    defer log.Println(time.Since(start))
    // ... do work ...
}
```

### Valid

```golang
import (
    "log"
    "time"
)

func example() {
    start := time.Now()
    // Good: wrap in a closure so time.Since is evaluated at defer time
    defer func() {
        log.Println(time.Since(start))
    }()
    // ... do work ...
}
```
