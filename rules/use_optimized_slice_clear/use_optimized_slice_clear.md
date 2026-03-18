---
title: useOptimizedSliceClear
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useOptimizedSliceClear`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useOptimizedSliceClear:
    # no additional options
```

## Details

Detects slice clear loops that can use compiler-optimized idiom. Go's compiler can optimize `for i := range s { s[i] = zero }` patterns into a `memclr` call. This checker identifies loops that clear slices and suggests using the optimized pattern.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
for i := 0; i < len(s); i++ {
    s[i] = 0
}
```

### Valid

```golang
for i := range s {
    s[i] = 0
}
```

```golang
// In Go 1.21+, you can use the clear builtin
clear(s)
```
