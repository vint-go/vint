---
title: noNilVariableReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noNilVariableReturn`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noNilVariableReturn:
    # no additional options
```

## Details

Detects return statements whose results evaluate to nil. This checker identifies conditional statements where a variable is compared to nil and then that same nil variable is returned within the if block. This suggests either an explicit `nil` return should replace the variable, or the comparison operator may contain a typo (`==` should be `!=`).

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f(x *int) *int {
    if x == nil {
        return x // returns nil, should just return nil explicitly or check != nil
    }
    return x
}
```

### Valid

```golang
func f(x *int) *int {
    if x == nil {
        return nil // explicit nil return
    }
    return x
}
```

```golang
func f(x *int) *int {
    if x != nil {
        return x // correct check before returning
    }
    return nil
}
```
