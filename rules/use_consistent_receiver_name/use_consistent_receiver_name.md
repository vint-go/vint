---
title: useConsistentReceiverName
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useConsistentReceiverName`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useConsistentReceiverName:
    # no additional options
```

## Details

Use consistent method receiver names.

All methods on a given type should use the same receiver name. Mixing different names (e.g., `s` in one method and `srv` in another) is inconsistent and confusing.

Source: https://staticcheck.dev/docs/checks/#ST1016

## Examples

### Invalid

```golang
package main

type Server struct{}

func (s *Server) Start() {}
func (srv *Server) Stop() {}
```

### Valid

```golang
package main

type Server struct{}

func (s *Server) Start() {}
func (s *Server) Stop() {}
```
