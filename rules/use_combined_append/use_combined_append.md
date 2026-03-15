---
title: useCombinedAppend
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useCombinedAppend`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useCombinedAppend:
    # no additional options
```

## Details

Detects consecutive append operations on the same slice that could be combined into a single call for better efficiency. When multiple sequential `append` calls target the same slice, they can be merged into one call with multiple arguments, which is both more readable and more performant.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
xs = append(xs, 1)
xs = append(xs, 2)
```

```golang
result = append(result, a)
result = append(result, b)
result = append(result, c)
```

### Valid

```golang
xs = append(xs, 1, 2)
```

```golang
result = append(result, a, b, c)
```
