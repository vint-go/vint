---
title: usePointerInSyncPool
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/usePointerInSyncPool`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/usePointerInSyncPool:
    # rule options here
```

## Details

Storing non-pointer values in `sync.Pool` allocates memory.

`sync.Pool` stores values as `interface{}`. Storing a non-pointer value causes an allocation because the value must be boxed into an interface. Store pointers instead to avoid this overhead.

Source: https://staticcheck.dev/docs/checks/#SA6002

## Examples

### Invalid

```golang
package main

import "sync"

var pool = sync.Pool{
    New: func() interface{} {
        // Allocates: storing a non-pointer value
        return make([]byte, 1024)
    },
}
```

### Valid

```golang
package main

import "sync"

var pool = sync.Pool{
    New: func() interface{} {
        b := make([]byte, 1024)
        return &b
    },
}
```
