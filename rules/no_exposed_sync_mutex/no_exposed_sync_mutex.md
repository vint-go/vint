---
title: noExposedSyncMutex
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noExposedSyncMutex`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noExposedSyncMutex:
    # no additional options
```

## Details

Detects exported methods from `sync.Mutex` and `sync.RWMutex` that are inadvertently promoted to the public API of a struct. When a struct embeds a mutex as an exported field or as an unnamed embedded type, the `Lock` and `Unlock` methods become part of the struct's public API, which is usually unintended. The mutex should be stored in an unexported field to keep synchronization as an implementation detail.

This rule is distinct from `noCopiedLock`, which detects copying of sync types by value. This rule focuses on API design by flagging public exposure of locking methods.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
type MyStruct struct {
    sync.Mutex // exported: Lock/Unlock become public methods
    Data string
}
```

### Valid

```golang
type MyStruct struct {
    mu   sync.Mutex // unexported: Lock/Unlock stay private
    Data string
}
```
