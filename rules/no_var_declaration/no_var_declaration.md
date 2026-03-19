---
title: noVarDeclaration
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noVarDeclaration`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noVarDeclaration:
    # no configuration options
```

## Details

This rule proposes simplifications of variable declarations. It detects two patterns:

1. **Zero value assignments**: When a variable is explicitly initialized with its zero value (e.g., `var x int = 0`), the assignment is redundant because Go automatically initializes variables to their zero values.

2. **Redundant type declarations**: When a variable is declared with an explicit type that matches the type of the right-hand side expression (e.g., `var x int = someIntFunc()`), the type annotation is redundant because Go can infer it from the right-hand side.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
var myInt int = 7 // should omit type int; it will be inferred
```

```golang
var mux *http.ServeMux = http.NewServeMux() // should omit type; it will be inferred
```

```golang
var myZeroInt int = 0 // should drop = 0; it is the zero value
```

```golang
var myZeroStr string = "" // should drop = ""; it is the zero value
```

```golang
var myZeroPtr *MyType = nil // should drop = nil; it is the zero value
```

### Valid

```golang
var x = 0 // no explicit type, nothing to simplify
```

```golang
var str fmt.Stringer // no initial value, nothing to simplify
```

```golang
var userID int64 = 1235 // LHS type differs from RHS default type
```

```golang
var out io.Writer = os.Stdout // LHS is interface, RHS is concrete type
```

```golang
var _ Server = (*serverImpl)(nil) // common idiom for interface satisfaction
```
