---
title: useOptimalOperandsOrder
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useOptimalOperandsOrder`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useOptimalOperandsOrder:
    # rule options here
```

## Details

Conditional expressions can be written to take advantage of short circuit evaluation and speed up their average evaluation time by forcing the evaluation of less time-consuming terms before more costly ones. This rule spots logical expressions where the order of evaluation of terms seems non-optimal.

Please notice that the confidence of this rule is low and it is up to the user to decide if the suggested rewrite of the expression keeps the semantics of the original one.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
if isGenerated(content) && !config.IgnoreGeneratedHeader {
    // function call on the left, simple field access on the right
}
```

```golang
if caller(x, y) || y {
    // function call on the left, simple variable on the right
}
```

### Valid

```golang
if !config.IgnoreGeneratedHeader && isGenerated(content) {
    // simple field access evaluated first, function call second
}
```

```golang
if y || caller(x, y) {
    // simple variable evaluated first, function call second
}
```
