---
title: noRedundantStringConcat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantStringConcat`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantStringConcat:
    # no additional options
```

## Details

Detects string concatenation operations that can be simplified. This checker identifies redundant string concatenation patterns where an empty string literal is concatenated with another value, and suggests simpler alternatives.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
s := "" + x // empty string concatenation
```

```golang
s := x + "" // empty string concatenation
```

### Valid

```golang
s := x
```
