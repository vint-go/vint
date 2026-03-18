---
title: useWriteByte
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useWriteByte`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useWriteByte:
    # no additional options
```

## Details

Detects `WriteRune` calls with a single-byte rune argument that could be replaced with `WriteByte`. When writing a single ASCII character, `WriteByte` is more efficient than `WriteRune` as it avoids UTF-8 encoding overhead.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
w.WriteRune('\n')
```

### Valid

```golang
w.WriteByte('\n')
```
