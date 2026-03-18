---
title: useTypeSwitchGuard
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTypeSwitchGuard`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTypeSwitchGuard:
    # no additional options
```

## Details

Detects type switches that can benefit from a type guard clause with a variable. This checker identifies type switch statements that repeatedly perform type assertions on the same value within case clauses. Using a type guard assignment (`switch v := v.(type)`) would be more efficient and cleaner than repeatedly asserting the type.

The checker skips cases with multiple types in a single case clause, as those result in `interface{}` type.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch v.(type) {
case *TypeA:
    handleA(v.(*TypeA)) // redundant type assertion
case *TypeB:
    handleB(v.(*TypeB)) // redundant type assertion
}
```

### Valid

```golang
switch v := v.(type) {
case *TypeA:
    handleA(v) // v is already *TypeA
case *TypeB:
    handleB(v) // v is already *TypeB
}
```
