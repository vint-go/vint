---
title: noEvalOrderDependency
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noEvalOrderDependency`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noEvalOrderDependency:
    # no additional options
```

## Details

Detects unwanted dependencies on evaluation order in return statements. This rule identifies problematic patterns where:

1. A return statement has multiple results, and one result is a variable that gets modified by a function call in another result.
2. A function call takes the address of a variable that is also being returned.
3. A method call uses a variable by pointer when that same variable appears in the return statement.

The evaluation order of return values is not guaranteed to be left-to-right in Go, so relying on side effects between return values leads to unpredictable behavior. Such function calls should be evaluated before the return statement.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f() (int, error) {
    x := 10
    return x, setX(&x) // x might be modified by setX
}
```

### Valid

```golang
func f() (int, error) {
    x := 10
    err := setX(&x)
    return x, err
}
```
