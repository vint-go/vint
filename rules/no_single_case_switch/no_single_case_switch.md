---
title: noSingleCaseSwitch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noSingleCaseSwitch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noSingleCaseSwitch:
    # no additional options
```

## Details

Detects switch statements that could be better written as if statements. This checker identifies two problematic patterns:

1. **Single case switches** -- Switch statements with exactly one case clause (not default), which should be converted to an `if` statement.
2. **Default-only switches** -- Switch statements containing only a default case, which is redundant.

The checker excludes cases where the case clause contains a `break` statement.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch x {
case 1:
    handleOne()
}
```

```golang
switch {
default:
    doSomething()
}
```

### Valid

```golang
if x == 1 {
    handleOne()
}
```

```golang
doSomething()
```

```golang
switch x {
case 1:
    handleOne()
case 2:
    handleTwo()
}
```
