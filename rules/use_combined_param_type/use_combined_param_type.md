---
title: useCombinedParamType
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCombinedParamType`
- This rule is not recommended (opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCombinedParamType:
    # no additional options
```

## Details

Detects if function parameters could be combined by type and suggests the way to do it. This checker identifies opportunities to consolidate function parameters that share the same type into a more concise declaration. The checker skips unnamed parameter lists and multi-line parameter declarations to avoid false positives.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func foo(a, b int, c, d int, e, f int, g int) {}
```

### Valid

```golang
func foo(a, b, c, d, e, f, g int) {}
```
