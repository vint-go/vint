---
title: noRedundantCanonicalHeaderKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantCanonicalHeaderKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantCanonicalHeaderKey:
    # rule options here
```

## Details

Redundant call to `net/http.CanonicalHeaderKey` in method call on `net/http.Header`.

Methods on `http.Header` already canonicalize the key internally. Calling `http.CanonicalHeaderKey` before passing the key is redundant.

Source: https://staticcheck.dev/docs/checks/#S1035

## Examples

### Invalid

```golang
package main

import "net/http"

func getHeader(h http.Header) string {
    return h.Get(http.CanonicalHeaderKey("content-type"))
}
```

### Valid

```golang
package main

import "net/http"

func getHeader(h http.Header) string {
    return h.Get("content-type")
}
```
