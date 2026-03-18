---
title: useElseIf
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useElseIf`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useElseIf:
    skipBalanced: true  # skip cases where if and else have single statements (default: true)
```

## Details

Detects `else` blocks that contain only a nested `if` statement, which can be simplified into an `else if` construct. Flattening the nesting improves readability and reduces unnecessary indentation. The `skipBalanced` parameter (default: true) optionally excludes cases where both the outer `if` and inner `if` contain only a single statement each.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
if a {
    doA()
} else {
    if b {
        doB()
    }
}
```

### Valid

```golang
if a {
    doA()
} else if b {
    doB()
}
```
