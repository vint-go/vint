---
title: noIgnoredQueryModification
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIgnoredQueryModification`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIgnoredQueryModification:
    # rule options here
```

## Details

`(*net/url.URL).Query` returns a copy, modifying it doesn't change the URL.

Calling `url.Query()` returns a new `url.Values` map. Modifying the returned map does not affect the original URL. You must call `url.RawQuery = values.Encode()` to apply changes.

Source: https://staticcheck.dev/docs/checks/#SA4027

## Examples

### Invalid

```golang
package main

import "net/url"

func addParam(u *url.URL) {
    // Modifying the copy, original URL is unchanged
    u.Query().Set("key", "value")
}
```

### Valid

```golang
package main

import "net/url"

func addParam(u *url.URL) {
    q := u.Query()
    q.Set("key", "value")
    u.RawQuery = q.Encode()
}
```
