---
title: noEmptyFallthrough
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noEmptyFallthrough`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noEmptyFallthrough:
    # no additional options
```

## Details

Detects `switch` case clauses that contain only a `fallthrough` statement with no other logic. Such cases can be simplified by combining the case values into a comma-separated list, which is more idiomatic in Go and easier to read.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
switch x {
case 1:
    fallthrough
case 2:
    handleOneOrTwo()
}
```

```golang
switch x {
case 3:
    fallthrough
default:
    handleDefault()
}
```

### Valid

```golang
switch x {
case 1, 2:
    handleOneOrTwo()
}
```

```golang
switch x {
default:
    handleDefault()
}
```
