---
title: noWaitGroupMisuse
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noWaitGroupMisuse`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noWaitGroupMisuse:
    # rule options here
```

## Details

Detects simple misuses of `sync.WaitGroup`. A common mistake is calling `wg.Add` inside a goroutine instead of before launching the goroutine. This creates a race condition where `wg.Wait` might return before all goroutines have started and called `wg.Add`.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/waitgroup

## Examples

### Invalid

```golang
import "sync"

func example() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        go func() {
            wg.Add(1) // Bad: Add called inside goroutine, race with Wait
            defer wg.Done()
            doWork()
        }()
    }
    wg.Wait()
}
```

### Valid

```golang
import "sync"

func example() {
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1) // Good: Add called before launching goroutine
        go func() {
            defer wg.Done()
            doWork()
        }()
    }
    wg.Wait()
}
```
