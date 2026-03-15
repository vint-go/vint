---
title: noUnnecessaryDeferLambda
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryDeferLambda`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryDeferLambda:
    # no additional options
```

## Details

Detects deferred function literals that can be simplified. When a deferred function literal only calls a single function, the wrapper is unnecessary and the function can be deferred directly. This makes the code cleaner and more readable.

This rule is distinct from `lint/correctness/noDeferTimeMisuse`, which detects cases where wrapping in a closure is actually required (e.g., `defer func() { log.Println(time.Since(start)) }()`). This rule only flags closures that serve no purpose.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
defer func() { f.Close() }()
```

### Valid

```golang
defer f.Close()
```
