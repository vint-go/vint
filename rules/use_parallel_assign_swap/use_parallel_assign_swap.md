---
title: useParallelAssignSwap
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useParallelAssignSwap`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useParallelAssignSwap:
    # no additional options
```

## Details

Detects value swapping code that does not use parallel assignment. Go supports parallel assignment which makes value swaps cleaner and more idiomatic. Using a temporary variable for swapping is unnecessary and adds visual noise.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
tmp := a
a = b
b = tmp
```

### Valid

```golang
a, b = b, a
```
