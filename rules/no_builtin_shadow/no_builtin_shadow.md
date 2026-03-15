---
title: noBuiltinShadow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noBuiltinShadow`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noBuiltinShadow:
    # no additional options
```

## Details

Detects when predeclared identifiers are shadowed in assignments. This checker identifies instances where built-in function or type names (such as `len`, `cap`, `append`, `make`, `new`, `error`, etc.) are reassigned as variable names, which shadows the original predeclared identifiers and makes them inaccessible within that scope.

This rule is distinct from `lint/suspicious/noVariableShadowing`, which detects general variable shadowing via `:=` in inner scopes. This rule specifically targets shadowing of Go's predeclared (builtin) identifiers in assignment statements.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Shadows the built-in len function
len := 10
```

```golang
// Shadows the built-in error type
error := fmt.Errorf("something failed")
```

### Valid

```golang
length := 10
```

```golang
err := fmt.Errorf("something failed")
```
