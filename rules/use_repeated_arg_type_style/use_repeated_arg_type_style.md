---
title: useRepeatedArgTypeStyle
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useRepeatedArgTypeStyle`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useRepeatedArgTypeStyle:
    arguments:
      - "short"  # or "full" or "any"
```

Or with separate configuration for arguments and return values:

```yaml title="vint.yaml"
settings:
  lint/style/useRepeatedArgTypeStyle:
    arguments:
      - funcArgStyle: "full"
        funcRetValStyle: "short"
```

## Details

This rule is designed to maintain consistency in the declaration of repeated argument and return value types in Go functions. It supports three styles: `any`, `short`, and `full`.

- `"any"`: No enforcement (default). Allows any form of type declaration.
- `"short"`: Encourages omitting repeated types for conciseness. Flags cases where consecutive parameters or return values share the same type but the type is written out explicitly for each one.
- `"full"`: Mandates explicitly stating the type for each argument and return value, even if they are repeated. Flags cases where multiple names share a single type declaration.

The rule can be configured as a single string (applied to both arguments and return values) or as a map with separate keys `funcArgStyle` and `funcRetValStyle` for independent control.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With "short" style:
func foo(a int, b int, c string) {}  // repeated argument type "int" can be omitted

// With "full" style:
func bar(a, b int, c string) {}  // argument types should not be omitted

// With "full" style for return values:
func baz() (a, b int, c string) { panic("implement me") }  // return types should not be omitted

// With "short" style for return values:
func qux() (a int, b int, c string) { panic("implement me") }  // repeated return type "int" can be omitted
```

### Valid

```golang
// With "short" style:
func foo(a, b int, c string) {}

// With "full" style:
func bar(a int, b int, c string) {}

// With "full" style for return values:
func baz() (a int, b int, c string) { panic("implement me") }

// With "short" style for return values:
func qux() (a, b int, c string) { panic("implement me") }
```
