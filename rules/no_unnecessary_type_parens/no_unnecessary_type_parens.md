---
title: noUnnecessaryTypeParens
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryTypeParens`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryTypeParens:
    # no additional options
```

## Details

Detects unneeded parentheses inside type expressions and suggests removing them. This checker identifies and flags unnecessary parentheses in Go type expressions such as array element types, pointer base types, type assertion types, function parameter and return types, map key and value types, and channel value types.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var x *(int)
```

```golang
var m map[(string)](int)
```

### Valid

```golang
var x *int
```

```golang
var m map[string]int
```

```golang
// Pointer dereference with parentheses is not flagged (not a type expression)
val := *(ptr)
```
