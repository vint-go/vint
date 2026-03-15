---
title: noBadLockPattern
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noBadLockPattern`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noBadLockPattern:
    # no additional options
```

## Details

Detects suspicious mutex lock/unlock operations. This checker identifies patterns where a mutex is locked or unlocked incorrectly, such as locking a mutex twice without an unlock in between, or calling Unlock immediately after Lock without doing any work in the critical section. These patterns typically indicate bugs.

Note: This rule is distinct from `lint/correctness/noCopiedLock`, which detects copying of sync types by value. This rule focuses on incorrect Lock/Unlock call sequences.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
mu.Lock()
mu.Lock() // double lock without unlock
```

```golang
mu.Lock()
mu.Unlock()
mu.Unlock() // double unlock
```

### Valid

```golang
mu.Lock()
defer mu.Unlock()
// critical section
```

```golang
mu.Lock()
doWork()
mu.Unlock()
```
