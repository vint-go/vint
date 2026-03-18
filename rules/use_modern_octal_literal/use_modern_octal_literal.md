---
title: useModernOctalLiteral
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useModernOctalLiteral`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useModernOctalLiteral:
    # no additional options
```

## Details

Detects old-style octal literals. This checker identifies integer literals using the legacy octal notation format (starting with `0` followed by digits) and suggests converting them to the modern `0o` prefix style introduced in Go 1.13. The checker only applies when the Go version is 1.13 or later.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
mode := 0755 // old-style octal literal
```

### Valid

```golang
mode := 0o755 // modern octal literal
```
