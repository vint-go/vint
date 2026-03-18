---
title: useSimplifiedSelector
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSimplifiedSelector`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSimplifiedSelector:
    # no additional options
```

## Details

Omit embedded fields from selector expression.

When accessing a field or method on an embedded type, you can omit the intermediate field name if there is no ambiguity.

Source: https://staticcheck.dev/docs/checks/#QF1008

## Examples

### Invalid

```golang
package main

import "sync"

type Server struct {
    sync.Mutex
}

func (s *Server) Process() {
    s.Mutex.Lock()
    defer s.Mutex.Unlock()
}
```

### Valid

```golang
package main

import "sync"

type Server struct {
    sync.Mutex
}

func (s *Server) Process() {
    s.Lock()
    defer s.Unlock()
}
```
