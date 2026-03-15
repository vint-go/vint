---
title: noCapitalizedLocal
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noCapitalizedLocal`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noCapitalizedLocal:
    paramsOnly: true  # whether to restrict checker to params only (default: true)
```

## Details

Detects capitalized names for local variables. In Go, capitalized names are used to denote exported identifiers. Using capitalized names for local variables, function parameters, or return values violates Go naming conventions and can cause confusion.

When `paramsOnly` is enabled (default), the checker only examines function parameters and return values. When disabled, it also checks local variable declarations.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f(IN int, OUT *int) (ERR error) {
    return nil
}
```

### Valid

```golang
func f(in int, out *int) (err error) {
    return nil
}
```
