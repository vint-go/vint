---
title: useSyncMapLoadAndDelete
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useSyncMapLoadAndDelete`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useSyncMapLoadAndDelete:
    # no additional options
```

## Details

Detects `sync.Map` load followed by delete that can be replaced with `LoadAndDelete`. When a `sync.Map.Load` call is immediately followed by a `sync.Map.Delete` on the same key, the atomic `LoadAndDelete` method (available since Go 1.15) should be used instead for correctness and efficiency.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
v, ok := m.Load(key)
m.Delete(key)
```

### Valid

```golang
v, ok := m.LoadAndDelete(key)
```
