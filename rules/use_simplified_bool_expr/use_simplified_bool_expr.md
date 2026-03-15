---
title: useSimplifiedBoolExpr
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/useSimplifiedBoolExpr`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/useSimplifiedBoolExpr:
    # no additional options
```

## Details

Detects bool expressions that can be simplified. This checker identifies and suggests simplifications for boolean expressions, including:

- **Double negation:** Removes redundant NOT operators (e.g., `!(!(x))` to `x`).
- **Negated equality:** Simplifies `!(x) == !(y)` to `(x) == (y)`.
- **Inverted comparisons:** Converts negated comparisons (e.g., `!(x >= y)` to `x < y`).
- **Combined checks:** Merges adjacent comparisons (e.g., `x > c || x == c` to `x >= c`).
- **Increment/decrement removal:** Simplifies off-by-one patterns (e.g., `x > y-1` to `x >= y`).
- **Range folding:** Combines range checks (e.g., `x >= c && x <= c` to `x == c`).

The checker avoids floating-point expressions to prevent type-related issues.

This rule is distinct from `lint/suspicious/noRedundantBoolCondition`, which detects redundant or duplicate conditions (e.g., `x == 1 || x == 1`). This rule focuses on algebraic simplification of boolean expressions.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Double negation
if !(!(ready)) {
    process()
}
```

```golang
// Can simplify negated comparison
if !(x >= y) {
    handle()
}
```

```golang
// Can combine comparisons
if x > 10 || x == 10 {
    handle()
}
```

### Valid

```golang
if ready {
    process()
}
```

```golang
if x < y {
    handle()
}
```

```golang
if x >= 10 {
    handle()
}
```
