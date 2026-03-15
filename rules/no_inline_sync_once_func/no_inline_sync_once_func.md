---
title: noInlineSyncOnceFunc
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInlineSyncOnceFunc`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInlineSyncOnceFunc:
    # no additional options
```

## Details

Detects bad usage of `sync.OnceFunc`. This checker identifies patterns where `sync.OnceFunc` is used incorrectly, such as calling it inline instead of storing the result. `sync.OnceFunc` returns a function that should be stored and called multiple times -- it ensures the wrapped function only executes once.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Calling sync.OnceFunc inline defeats its purpose
sync.OnceFunc(func() {
    initialize()
})()
```

### Valid

```golang
var initOnce = sync.OnceFunc(func() {
    initialize()
})

func doWork() {
    initOnce() // safe to call multiple times
}
```
