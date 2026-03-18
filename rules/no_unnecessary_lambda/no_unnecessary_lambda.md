---
title: noUnnecessaryLambda
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryLambda`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryLambda:
    # no additional options
```

## Details

Detects function literals that can be simplified. This checker identifies function literals that unnecessarily wrap a single function call and can be replaced with a direct reference to that function. The checker validates that all parameters are passed directly to the wrapped function in the same order and the function types match exactly.

This rule is distinct from `noUnnecessaryDeferLambda`, which specifically targets unnecessary closures in defer statements. This rule covers the general case of any unnecessary function literal wrapper.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
sort.Slice(xs, func(i, j int) bool {
    return less(i, j)
})
```

```golang
go func() { doWork() }()
```

### Valid

```golang
sort.Slice(xs, less)
```

```golang
go doWork()
```
