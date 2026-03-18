---
title: useFilepathJoin
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useFilepathJoin`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useFilepathJoin:
    # no additional options
```

## Details

Detects path concatenation that can be replaced with `filepath.Join`. Manual path concatenation using `+` and `/` is error-prone, especially across operating systems. Using `filepath.Join` handles path separator differences automatically.

This rule is distinct from `noSeparatorInFilepathJoin`, which detects separator characters inside arguments already passed to `filepath.Join`. This rule instead detects cases where `filepath.Join` is not used at all.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
path := dir + "/" + filename
```

### Valid

```golang
path := filepath.Join(dir, filename)
```
