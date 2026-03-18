---
title: useDecodeRune
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useDecodeRune`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useDecodeRune:
    # no additional options
```

## Details

Detects expressions like `[]rune(s)[0]` that perform an unnecessary full string-to-rune-slice conversion when only the first rune is needed. Using `utf8.DecodeRuneInString(s)` is more efficient as it avoids allocating a rune slice for the entire string.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
firstRune := []rune(s)[0]
```

### Valid

```golang
firstRune, _ := utf8.DecodeRuneInString(s)
```
