---
title: noExcessiveArguments
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noExcessiveArguments`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noExcessiveArguments:
    arguments: [4]
```

## Details

Warns when a function receives more parameters than the maximum set by the rule's configuration.
Enforcing a maximum number of parameters helps to keep the code readable and maintainable.

Functions with too many parameters are hard to understand and maintain. They often indicate
that the function is doing too much or that related parameters should be grouped into a struct.

The default limit is 8 parameters.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// with limit set to 4
func processOrder(id int, name string, price float64, quantity int, discount float64) {
    // too many parameters
}
```

### Valid

```golang
type OrderParams struct {
    ID       int
    Name     string
    Price    float64
    Quantity int
    Discount float64
}

func processOrder(params OrderParams) {
    // parameters grouped into a struct
}
```
