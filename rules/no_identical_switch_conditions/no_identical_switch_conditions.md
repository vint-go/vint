---
title: noIdenticalSwitchConditions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIdenticalSwitchConditions`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIdenticalSwitchConditions:
    # no configuration options
```

## Details

A `switch` statement with cases with the same condition can lead to unreachable code and is a potential source of bugs while making the code harder to read and maintain.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
switch {
case a > 0, a < 0:
case a == 0:
case a < 0: // duplicate condition: a < 0 already appears above
default:
}
```

```golang
switch {
case lnOpts.IsSocketOpts():
    fallthrough
case lnOpts.IsTimeout(), lnOpts.IsSocketOpts(): // duplicate: lnOpts.IsSocketOpts()
    // ...
case lnOpts.IsTimeout(): // duplicate: lnOpts.IsTimeout()
    // ...
default:
}
```

### Valid

```golang
switch {
case a > 0:
case a == 0:
case a < 0:
default:
}
```

```golang
switch expression {
case value1:
case value2:
default:
}
```
