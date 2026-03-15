---
title: noUnreachableTypeCase
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnreachableTypeCase`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnreachableTypeCase:
    # no additional options
```

## Details

Detects erroneous case order inside switch statements. This checker identifies problematic case ordering in type switch statements where a concrete type case appears after an interface type case that would match it, making the concrete case unreachable.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// ast.Expr (interface) is matched before *ast.BasicLit (concrete),
// making the concrete case unreachable
switch v := x.(type) {
case ast.Expr:
    // handles all expressions
case *ast.BasicLit:
    // unreachable: already caught by ast.Expr
}
```

### Valid

```golang
// Concrete type is matched before interface type
switch v := x.(type) {
case *ast.BasicLit:
    // handles basic literals
case ast.Expr:
    // handles other expressions
}
```
