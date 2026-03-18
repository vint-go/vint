---
title: useRegexpMustCompile
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useRegexpMustCompile`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useRegexpMustCompile:
    # no additional options
```

## Details

Detects `regexp.Compile*` calls that can be replaced with `regexp.MustCompile*`. When a regular expression is compiled with a constant pattern at package level or in an `init` function, `MustCompile` is preferred because it panics on invalid patterns, catching errors at startup rather than at runtime.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
var re, _ = regexp.Compile(`\d+`)
```

### Valid

```golang
var re = regexp.MustCompile(`\d+`)
```
