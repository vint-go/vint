---
title: noDuplicateSubExpression
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateSubExpression`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateSubExpression:
    # no additional options
```

## Details

Detects suspicious duplicated sub-expressions. This checker identifies binary expressions where both operands are identical, which typically indicates a logic error (often a copy-paste mistake). It flags cases like `x < x` or `a && a` where the same subexpression appears on both sides of an operator.

The checker skips floating-point comparisons in certain operations (like equality checks) to avoid false positives related to NaN behavior, and only warns when expressions are side-effect-free.

This rule is broader than `lint/suspicious/noRedundantBoolCondition`, which specifically detects duplicate conditions in `&&` and `||` expressions. This rule covers all binary operators including comparison, arithmetic, and logical operators.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Same subexpression on both sides
if xs[i].v < xs[i].v {
    // always false
}
```

```golang
// Identical operands in logical expression
if a || a {
    doSomething()
}
```

### Valid

```golang
if xs[i].v < xs[j].v {
    // different indices
}
```

```golang
if a || b {
    doSomething()
}
```
