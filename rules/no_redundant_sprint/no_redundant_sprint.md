---
title: noRedundantSprint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantSprint`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantSprint:
    # no additional options
```

## Details

Detects redundant `fmt.Sprint` calls. When a value is already a string, calling `fmt.Sprint` or `fmt.Sprintf("%s", ...)` on it is unnecessary and adds overhead. The checker identifies these cases and suggests using the string value directly.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
s := "hello"
result := fmt.Sprint(s) // redundant: s is already a string
```

```golang
result := fmt.Sprintf("%s", str) // redundant conversion
```

### Valid

```golang
s := "hello"
result := s
```

```golang
result := fmt.Sprintf("prefix: %s", str) // not redundant, has additional formatting
```
