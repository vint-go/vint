---
title: noImpossibleCondition
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noImpossibleCondition`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noImpossibleCondition:
    # no additional options
```

## Details

Detects suspicious condition expressions that are likely logic errors. This checker identifies several problematic patterns:

1. **Loop condition errors** -- Inverted comparison operators in for loops (e.g., `i > n` instead of `i < n`).
2. **Impossible conditions** -- Contradictory boolean expressions like `x < a && x > b` where `a < b` (always false).
3. **Suspicious equality chains** -- Patterns like `x == a && x == b` where a variable is compared for equality to different values simultaneously.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Inverted loop condition
for i := 0; i > n; i++ {
    // body never executes
}
```

```golang
// Impossible condition: x cannot be both less than 5 and greater than 10
if x < 5 && x > 10 {
    // unreachable
}
```

```golang
// x cannot equal two different values at the same time
if x == 1 && x == 2 {
    // unreachable
}
```

### Valid

```golang
for i := 0; i < n; i++ {
    // correct loop condition
}
```

```golang
if x > 5 && x < 10 {
    // valid range check
}
```
