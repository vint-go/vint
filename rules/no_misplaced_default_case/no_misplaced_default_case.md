---
title: noMisplacedDefaultCase
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMisplacedDefaultCase`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMisplacedDefaultCase:
    # no additional options
```

## Details

Detects when default case in switch isn't on first or last position. For readability and convention, the `default` case in a switch statement should be placed at either the first or last position. Having it in the middle makes the switch harder to read and understand.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch x {
case 1:
    handleOne()
default:
    handleDefault()
case 2:
    handleTwo()
}
```

### Valid

```golang
switch x {
case 1:
    handleOne()
case 2:
    handleTwo()
default:
    handleDefault()
}
```

```golang
switch x {
default:
    handleDefault()
case 1:
    handleOne()
case 2:
    handleTwo()
}
```
