---
title: useTimeMethod
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTimeMethod`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTimeMethod:
    # no additional options
```

## Details

Detects manual conversions to milli- or microseconds. Instead of manually multiplying or dividing time durations to convert units, the `time` package's built-in methods and constants should be used for clarity and correctness.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
d := t.UnixNano() / 1000000 // manual conversion to milliseconds
```

### Valid

```golang
d := t.UnixMilli()
```

```golang
d := t.UnixMicro()
```
