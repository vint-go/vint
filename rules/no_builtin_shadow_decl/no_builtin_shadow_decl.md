---
title: noBuiltinShadowDecl
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noBuiltinShadowDecl`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noBuiltinShadowDecl:
    # no additional options
```

## Details

Detects top-level declarations that shadow the predeclared identifiers. This checker identifies when developers create top-level declarations (functions, types, or values) using names that conflict with Go's built-in identifiers. Methods are excluded since they can safely shadow names.

This rule is distinct from `lint/suspicious/noVariableShadowing`, which detects general variable shadowing via `:=` in inner scopes, and from `lint/suspicious/noBuiltinShadow`, which detects builtin shadowing in assignments. This rule specifically targets top-level (package-level) declarations that shadow predeclared identifiers.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Shadows the built-in int type
type int struct{}
```

```golang
// Shadows the built-in error type
var error = "something"
```

### Valid

```golang
type myInt struct{}
```

```golang
var errMsg = "something"
```
