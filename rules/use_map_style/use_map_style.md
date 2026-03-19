---
title: useMapStyle
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useMapStyle`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useMapStyle:
    arguments:
      - "make"  # or "literal" or "any"
```

## Details

This rule enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization. It does not affect `make(map[type]type, size)` constructions as well as `map[type]type{k1: v1}`.

The configuration option is a string specifying the enforced style for map initialization:

- `"any"`: No enforcement (default).
- `"make"`: Enforces the usage of `make(map[type]type)`.
- `"literal"`: Enforces the usage of `map[type]type{}`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With "make" style enforced:
m := map[string]string{}  // should use make(map[string]string)

// With "literal" style enforced:
m := make(map[string]string)  // should use map[string]string{}
```

### Valid

```golang
// With "make" style enforced:
m := make(map[string]string)
m2 := make(map[string]string, 10)  // sized make is always allowed
m3 := map[string]string{"k1": "v1"}  // non-empty literals are always allowed

// With "literal" style enforced:
m := map[string]string{}
m2 := make(map[string]string, 10)  // sized make is always allowed
m3 := map[string]string{"k1": "v1"}  // non-empty literals are always allowed
```
