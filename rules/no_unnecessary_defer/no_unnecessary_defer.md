---
title: noUnnecessaryDefer
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryDefer`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryDefer:
    # no additional options
```

## Details

Detects redundantly deferred calls. This checker identifies `defer` statements that are placed immediately before a function's return statement and serve no purpose. The deferred call could simply be executed directly without the `defer` mechanism, since the function is about to return anyway.

This rule is distinct from `lint/correctness/noDeferTimeMisuse`, which detects `defer` with `time.Since` that evaluates arguments immediately. It is also distinct from `lint/style/noUnnecessaryDeferLambda`, which detects unnecessary closures wrapping deferred function calls. This rule specifically targets defer statements that are the last statement before a return, making the defer mechanism redundant.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f() error {
    mu.Lock()
    // ... do work ...
    defer mu.Unlock() // unnecessary: right before return
    return nil
}
```

### Valid

```golang
func f() error {
    mu.Lock()
    defer mu.Unlock() // useful: there is code after this
    // ... do work ...
    return doSomething()
}
```

```golang
func f() error {
    mu.Lock()
    // ... do work ...
    mu.Unlock() // called directly instead of deferred
    return nil
}
```
