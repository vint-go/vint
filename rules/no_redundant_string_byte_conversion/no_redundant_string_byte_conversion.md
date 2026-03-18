---
title: noRedundantStringByteConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noRedundantStringByteConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noRedundantStringByteConversion:
    # no additional options
```

## Details

Detects redundant conversions between string and `[]byte`. When a function accepts `[]byte` and there is an equivalent function that accepts `string` (or vice versa), a conversion can be avoided. This checker identifies unnecessary `string(b)` and `[]byte(s)` conversions that add allocation overhead.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Unnecessary conversion from string to []byte
result := bytes.Contains([]byte(s), []byte("pattern"))
```

### Valid

```golang
// Use strings package directly
result := strings.Contains(s, "pattern")
```
