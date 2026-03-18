---
title: noRedundantSwitchTrue
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantSwitchTrue`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantSwitchTrue:
    # no additional options
```

## Details

Detects switch-over-bool statements that use explicit `true` tag value. In Go, `switch { ... }` is equivalent to `switch true { ... }`, so the explicit `true` is redundant.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch true {
case x > 0:
    handlePositive()
case x < 0:
    handleNegative()
}
```

### Valid

```golang
switch {
case x > 0:
    handlePositive()
case x < 0:
    handleNegative()
}
```
