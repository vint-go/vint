---
title: useSliceStyle
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSliceStyle`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSliceStyle:
    arguments:
      - "make"  # or "literal", "nil", or "any"
```

## Details

This rule enforces consistent usage of `make([]type, 0)`, `[]type{}`, or `var []type` for slice initialization. It does not affect `make([]type, non_zero_len, or_non_zero_cap)` constructions as well as `[]type{v1}`. Nil slices are always permitted.

The configuration option is a string specifying the enforced style for slice initialization:

- `"any"`: No enforcement (default).
- `"make"`: Enforces the usage of `make([]type, 0)`.
- `"literal"`: Enforces the usage of `[]type{}`.
- `"nil"`: Enforces the usage of `var []type` (nil slice declaration).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With "make" style enforced:
m := []string{}  // should use make([]string, 0)

// With "literal" style enforced:
m := make([]string, 0)  // should use []string{}

// With "nil" style enforced:
m := []string{}          // should use var m []string
m := make([]string, 0)  // should use var m []string
```

### Valid

```golang
// With "make" style enforced:
m := make([]string, 0)
m2 := make([]string, 10)        // non-zero size is always allowed
m3 := []string{"v1", "v2"}      // non-empty literals are always allowed
var m4 []string                  // nil declaration is always allowed

// With "literal" style enforced:
m := []string{}
m2 := make([]string, 10)        // non-zero size is always allowed
m3 := []string{"v1", "v2"}      // non-empty literals are always allowed

// With "nil" style enforced:
var m []string
m2 := make([]string, 10)        // non-zero size is always allowed
m3 := []string{"v1", "v2"}      // non-empty literals are always allowed
```
