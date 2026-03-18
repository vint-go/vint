---
title: useTypeSwitchChain
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeSwitchChain`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeSwitchChain:
    # no additional options
```

## Details

Detects repeated type assertions and suggests replacing them with type switch statements. This checker identifies chains of if-else statements that perform consecutive type assertions on the same variable. When it finds 2 or more type assertions in sequence, it recommends refactoring the code to use a type switch, which is more idiomatic Go.

This rule is distinct from `noUncheckedTypeAssertion`, which detects type assertions without the comma-ok idiom. This rule focuses on the stylistic improvement of replacing if-else chains of type assertions with a type switch statement.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if v, ok := x.(*TypeA); ok {
    handleA(v)
} else if v, ok := x.(*TypeB); ok {
    handleB(v)
} else if v, ok := x.(*TypeC); ok {
    handleC(v)
}
```

### Valid

```golang
switch v := x.(type) {
case *TypeA:
    handleA(v)
case *TypeB:
    handleB(v)
case *TypeC:
    handleC(v)
}
```
