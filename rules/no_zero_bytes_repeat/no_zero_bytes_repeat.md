---
title: noZeroBytesRepeat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noZeroBytesRepeat`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noZeroBytesRepeat:
    # no additional options
```

## Details

Detects `bytes.Repeat` calls with a count of 0. Calling `bytes.Repeat` with zero count always returns an empty byte slice, which is likely a bug or can be simplified. If an empty byte slice is intended, it should be created directly as `[]byte{}` instead.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
padding := bytes.Repeat([]byte(" "), 0) // always returns empty slice
```

### Valid

```golang
padding := bytes.Repeat([]byte(" "), n) // n > 0
```

```golang
// If you need an empty slice, create one directly
padding := []byte{}
```
