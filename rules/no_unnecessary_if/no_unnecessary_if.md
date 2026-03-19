---
title: noUnnecessaryIf
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryIf`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryIf:
    # no configuration options
```

## Details

Detects unnecessary `if-else` statements that return or assign a boolean value
based on a condition and suggests a simplified, direct return or assignment.
The `if-else` block is redundant because the condition itself is already a boolean expression.
The simplified version is immediately clearer, more idiomatic, and reduces cognitive load for the reader.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
if y <= 0 {
  z = true
} else {
  z = false
}
```

```golang
if x > 10 {
  return false
} else {
  return true
}
```

### Valid

```golang
z = y <= 0
```

```golang
return x <= 10
```
