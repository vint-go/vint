---
title: noIdenticalBranches
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIdenticalBranches`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIdenticalBranches:
    # no options
```

## Details

An `if-then-else` conditional with identical implementations in both branches is an error. When both branches execute the same code, the condition is meaningless and the duplication likely indicates a copy-paste mistake or unfinished logic.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
if condition {
    doSomething()
} else {
    doSomething()
}
```

```golang
if x > 0 {
    // empty
} else {
    // empty
}
```

### Valid

```golang
if condition {
    doSomething()
} else {
    doSomethingElse()
}
```

```golang
if condition {
    doSomething()
}
```
