---
title: noOverlappingEncoderSlice
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noOverlappingEncoderSlice`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noOverlappingEncoderSlice:
    # rule options here
```

## Details

Overlapping byte slices passed to an encoder.

Some encoding functions, such as those in `encoding/base64` and `encoding/hex`, do not support overlapping source and destination byte slices. Using overlapping slices may produce incorrect results.

Source: https://staticcheck.dev/docs/checks/#SA1031

## Examples

### Invalid

```golang
package main

import "encoding/hex"

func main() {
    buf := make([]byte, 100)
    // Overlapping source and destination
    hex.Encode(buf, buf[:50])
}
```

### Valid

```golang
package main

import "encoding/hex"

func main() {
    src := []byte("hello")
    dst := make([]byte, hex.EncodedLen(len(src)))
    // Non-overlapping source and destination
    hex.Encode(dst, src)
}
```
