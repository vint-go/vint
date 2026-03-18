---
title: noSeparatorInFilepathJoin
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noSeparatorInFilepathJoin`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noSeparatorInFilepathJoin:
    # no additional options
```

## Details

Detects problems in `filepath.Join()` function calls. This checker identifies string literal arguments to `filepath.Join()` that contain path separators (forward slashes or backslashes). Since `filepath.Join()` handles path separators automatically, including them in arguments is redundant and potentially error-prone.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
path := filepath.Join("dir/", filename)
```

```golang
path := filepath.Join("base\\sub", filename)
```

### Valid

```golang
path := filepath.Join("dir", filename)
```

```golang
path := filepath.Join("base", "sub", filename)
```
