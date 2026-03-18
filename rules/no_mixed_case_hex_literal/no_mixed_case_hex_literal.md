---
title: noMixedCaseHexLiteral
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noMixedCaseHexLiteral`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noMixedCaseHexLiteral:
    # no additional options
```

## Details

Detects hex literals that have mixed case letter digits. This checker identifies hex literals using uppercase `0X` prefix (suggesting lowercase `0x`) and mixed casing in hex digit letters. Consistent casing in hex literals improves readability.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
x := 0xfF
```

```golang
y := 0XFF
```

### Valid

```golang
x := 0xff
```

```golang
y := 0xFF
```
