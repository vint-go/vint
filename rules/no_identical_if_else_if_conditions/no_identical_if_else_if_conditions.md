---
title: noIdenticalIfElseIfConditions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIdenticalIfElseIfConditions`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIdenticalIfElseIfConditions:
    # no configuration options
```

## Details

An `if ... else if` chain with identical conditions can lead to unreachable code and is a potential source of bugs while making the code harder to read and maintain.

When two branches in an `if ... else if` chain share the same condition, the second branch can never be reached because the first matching condition will always execute. This typically indicates a copy-paste mistake where the developer forgot to update the condition.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// The second condition `a > 0` is identical to the first, making it unreachable.
if a > 0 {
    print("something")
} else if a < 0 {
    print("something else")
} else if a > 0 {
    println()
}
```

```golang
// Multiple identical conditions in the same chain.
if a > 0 {
    print("something")
} else if a == 0 {
    print("zero")
} else if a > 0 {
    println()
} else if a == 0 {
    print("other")
}
```

### Valid

```golang
// Each condition is unique.
if a > 0 {
    print("positive")
} else if a < 0 {
    print("negative")
} else if a == 0 {
    print("zero")
}
```

```golang
// Conditions with initializers are not checked to avoid false positives.
if x := compute(); x > 0 {
    print("positive")
} else if x := recompute(); x > 0 {
    print("recomputed positive")
}
```
