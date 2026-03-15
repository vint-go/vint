---
title: noLostCancel
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noLostCancel`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noLostCancel:
    # rule options here
```

## Details

Checks for failure to call a context cancellation function. When you create a derived context using `context.WithCancel`, `context.WithTimeout`, or `context.WithDeadline`, the returned cancel function must be called on all execution paths to release resources associated with the context. Failing to do so can cause resource leaks, such as goroutine leaks.

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/lostcancel

## Examples

### Invalid

```golang
import "context"

func example(ctx context.Context) {
    // Bad: cancel function is never called
    ctx, _ = context.WithCancel(ctx)
    doWork(ctx)
}
```

```golang
import "context"

func example(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, time.Second)
    // Bad: cancel is not called if doWork returns an error
    result, err := doWork(ctx)
    if err != nil {
        return err // cancel is never called on this path
    }
    cancel()
    return nil
}
```

### Valid

```golang
import "context"

func example(ctx context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel() // Good: always called
    doWork(ctx)
}
```

```golang
import "context"

func example(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, time.Second)
    defer cancel() // Good: called on all paths
    result, err := doWork(ctx)
    if err != nil {
        return err
    }
    return nil
}
```
