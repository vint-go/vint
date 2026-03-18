---
title: useShortVarDecl
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useShortVarDecl`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useShortVarDecl:
    # no additional options
```

## Details

Detects suspicious/confusing re-assignments. This checker identifies cases where a variable is reassigned within an if statement's initialization block when a short variable declaration (`:=`) would be more appropriate. It flags patterns where an if statement's init section reassigns an existing variable, the condition checks if that variable is not nil, and the statement body returns that same variable.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var err error
if err = doSomething(); err != nil {
    return err
}
```

### Valid

```golang
if err := doSomething(); err != nil {
    return err
}
```
