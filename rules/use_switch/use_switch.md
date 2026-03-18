---
title: useSwitch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSwitch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSwitch:
    minThreshold: 2  # minimum number of if-else branches to trigger (default: 2)
```

## Details

Detects repeated if-else statements and recommends converting them to switch statements. When a chain of if-else-if blocks meets or exceeds the configurable minimum threshold, the checker suggests rewriting the conditional logic using a switch statement instead, which is more idiomatic and readable in Go.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if x == 1 {
    handleOne()
} else if x == 2 {
    handleTwo()
} else if x == 3 {
    handleThree()
}
```

### Valid

```golang
switch x {
case 1:
    handleOne()
case 2:
    handleTwo()
case 3:
    handleThree()
}
```
