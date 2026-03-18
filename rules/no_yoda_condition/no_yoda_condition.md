---
title: noYodaCondition
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noYodaCondition`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noYodaCondition:
    # no additional options
```

## Details

Detects Yoda style expressions and suggests reordering them. Yoda conditions place the constant on the left side of the comparison operator, which is considered non-idiomatic in Go. This checker suggests swapping the operands for better readability.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if nil != err {
    return err
}
```

```golang
if 42 == x {
    doSomething()
}
```

### Valid

```golang
if err != nil {
    return err
}
```

```golang
if x == 42 {
    doSomething()
}
```
