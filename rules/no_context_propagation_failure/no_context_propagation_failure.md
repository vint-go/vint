---
title: noContextPropagationFailure
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noContextPropagationFailure`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noContextPropagationFailure:
    # rule options here
```

## Details

Detects context propagation failures that can lead to goroutine or resource leaks.

This rule uses SSA analysis to identify patterns where a context is not properly propagated to goroutines or long-running operations. When a parent context is cancelled but child goroutines or operations do not receive the cancellation signal, those goroutines continue to run and consume resources indefinitely, leading to goroutine leaks and resource exhaustion.

Proper context propagation is essential in Go's concurrency model. Every goroutine that performs work on behalf of a request should receive and respect the request's context, enabling graceful cancellation and timeout propagation throughout the call chain.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func handler(ctx context.Context) {
    // Goroutine does not receive the parent context
    go func() {
        // This goroutine will not be cancelled when ctx is cancelled
        result := longRunningOperation()
        processResult(result)
    }()
}
```

```golang
func fetchAll(ctx context.Context, urls []string) {
    for _, url := range urls {
        // Context not passed to HTTP request
        go func(u string) {
            resp, _ := http.Get(u) // Should use http.NewRequestWithContext
            defer resp.Body.Close()
        }(url)
    }
}
```

### Valid

```golang
func handler(ctx context.Context) {
    // Goroutine receives and respects the parent context
    go func(ctx context.Context) {
        select {
        case <-ctx.Done():
            return
        case result := <-doWork(ctx):
            processResult(result)
        }
    }(ctx)
}
```

```golang
func fetchAll(ctx context.Context, urls []string) {
    for _, url := range urls {
        go func(ctx context.Context, u string) {
            req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
            resp, _ := http.DefaultClient.Do(req)
            defer resp.Body.Close()
        }(ctx, url)
    }
}
```
