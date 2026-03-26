---
title: useIdiomaticReceiverName
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIdiomaticReceiverName`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIdiomaticReceiverName:
    # no additional options
```

## Details

Poorly chosen receiver name.

Receiver names should not be generic names like `this` or `self`, and should not be an underscore `_`.

Source: https://staticcheck.dev/docs/checks/#ST1006

## Examples

### Invalid

```golang
package main

type Server struct{}

func (this *Server) Start() {}
func (self *Server) Stop() {}
```

### Valid

```golang
package main

type Server struct{}

func (s *Server) Start() {}
func (s *Server) Stop() {}
```
