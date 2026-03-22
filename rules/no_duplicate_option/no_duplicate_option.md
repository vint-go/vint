---
title: noDuplicateOption
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateOption`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateOption:
    # no additional options
```

## Details

Detects duplicated option function arguments in variadic function calls. This checker identifies when the same option function is passed multiple times as arguments to a variadic function call. The second occurrence is flagged as redundant.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// withWidth is passed twice
widget := NewWidget(withWidth(10), withHeight(20), withWidth(30))
```

### Valid

```golang
widget := NewWidget(withWidth(10), withHeight(20))
```
