---
title: useBufferStringOrBytes
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useBufferStringOrBytes`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useBufferStringOrBytes:
    # rule options here
```

## Details

Use `bytes.Buffer.String` or `bytes.Buffer.Bytes`.

Converting a `bytes.Buffer` to string via `string(buf.Bytes())` can be simplified to `buf.String()`.

Source: https://staticcheck.dev/docs/checks/#S1030

## Examples

### Invalid

```golang
package main

import "bytes"

func process(buf bytes.Buffer) string {
    return string(buf.Bytes())
}
```

### Valid

```golang
package main

import "bytes"

func process(buf bytes.Buffer) string {
    return buf.String()
}
```
