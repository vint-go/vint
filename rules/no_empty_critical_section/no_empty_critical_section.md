---
title: noEmptyCriticalSection
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noEmptyCriticalSection`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noEmptyCriticalSection:
    # rule options here
```

## Details

Empty critical section, did you mean to `defer` the unlock?

Locking a mutex and immediately unlocking it without doing any work in between is pointless and likely a mistake. This usually happens when the `Unlock` was meant to be deferred.

Source: https://staticcheck.dev/docs/checks/#SA2001

## Examples

### Invalid

```golang
package main

import "sync"

var mu sync.Mutex

func process() {
    mu.Lock()
    mu.Unlock()
    // Work done outside the critical section
    // is not protected by the mutex
}
```

### Valid

```golang
package main

import "sync"

var mu sync.Mutex

func process() {
    mu.Lock()
    defer mu.Unlock()
    // Work is protected by the mutex
}
```
