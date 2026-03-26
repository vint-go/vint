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

Detects path concatenation using `string(os.PathSeparator)` that can be replaced with `filepath.Join`. Using `filepath.Join` handles path separator differences automatically and is more readable.

This rule matches the specific pattern `x + string(os.PathSeparator) + y`, matching the behavior of go-critic's `preferFilepathJoin` checker. It does **not** flag general string concatenation with literal `/` or `\` characters (e.g., URLs, MQTT topics, HTTP paths).

This rule is distinct from `noSeparatorInFilepathJoin`, which detects separator characters inside arguments already passed to `filepath.Join`. This rule instead detects cases where `filepath.Join` is not used at all.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
path := dir + string(os.PathSeparator) + filename
```

### Valid

```golang
path := filepath.Join(dir, filename)

// These are NOT flagged (literal slashes, not os.PathSeparator):
url := "https://example.com/" + path
topic := "devices/" + device + "/status"
```
