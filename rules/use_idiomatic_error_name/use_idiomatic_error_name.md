---
title: useIdiomaticErrorName
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIdiomaticErrorName`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIdiomaticErrorName:
    # rule options here
```

## Details

Poorly chosen name for error variable.

Error variables should be named `err` or have a prefix of `err` or `Err` (for exported variables). Sentinel errors (package-level error variables) should be named `ErrFoo`, not `ErrorFoo` or `FooError`.

Source: https://staticcheck.dev/docs/checks/#ST1012

## Examples

### Invalid

```golang
package main

import "errors"

var NotFoundError = errors.New("not found")
```

### Valid

```golang
package main

import "errors"

var ErrNotFound = errors.New("not found")
```
