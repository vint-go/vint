---
title: noNonCanonicalHeaderKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNonCanonicalHeaderKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNonCanonicalHeaderKey:
    # rule options here
```

## Details

Non-canonical key in `http.Header` map.

HTTP header keys are canonicalized by `net/http` using `textproto.CanonicalMIMEHeaderKey`. When accessing headers directly via the map, you must use the canonical form. For example, `"content-type"` should be `"Content-Type"`. Alternatively, use `http.Header.Get()` which handles canonicalization automatically.

Source: https://staticcheck.dev/docs/checks/#SA1008

## Examples

### Invalid

```golang
package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    // Non-canonical key - will not match
    ct := r.Header["content-type"]
    _ = ct
}
```

### Valid

```golang
package main

import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    // Using Get() which handles canonicalization
    ct := r.Header.Get("content-type")
    _ = ct

    // Or using canonical form directly
    ct2 := r.Header["Content-Type"]
    _ = ct2
}
```
