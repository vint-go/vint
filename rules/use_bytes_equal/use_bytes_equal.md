---
title: useBytesEqual
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useBytesEqual`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useBytesEqual:
    # rule options here
```

## Details

Replace call to `bytes.Compare` with `bytes.Equal`.

Using `bytes.Compare(a, b) == 0` to check for equality can be simplified to `bytes.Equal(a, b)`.

Source: https://staticcheck.dev/docs/checks/#S1004

## Examples

### Invalid

```golang
package main

import "bytes"

func isEqual(a, b []byte) bool {
    return bytes.Compare(a, b) == 0
}
```

### Valid

```golang
package main

import "bytes"

func isEqual(a, b []byte) bool {
    return bytes.Equal(a, b)
}
```
