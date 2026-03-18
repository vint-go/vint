---
title: useConvenienceFunc
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useConvenienceFunc`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useConvenienceFunc:
    # no additional options
```

## Details

Detects function calls that can be replaced with convenience wrappers from the Go standard library. The standard library provides many convenience functions that wrap common patterns. This checker identifies cases where using the wrapper function would be cleaner and more idiomatic.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
strings.SplitN(s, sep, -1) // equivalent to strings.Split
```

### Valid

```golang
strings.Split(s, sep)
```
