---
title: noInvalidBinaryArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidBinaryArg`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidBinaryArg:
    # rule options here
```

## Details

Unsupported argument to functions in `encoding/binary`.

This check verifies that values passed to `binary.Read`, `binary.Write`, and `binary.Size` are of types supported by the `encoding/binary` package. The package only supports fixed-size types (e.g., `int32`, `float64`) and structs/slices of such types. Variable-size types like `int`, `string`, or `map` are not supported.

Source: https://staticcheck.dev/docs/checks/#SA1003

## Examples

### Invalid

```golang
package main

import (
    "bytes"
    "encoding/binary"
)

func main() {
    var buf bytes.Buffer
    // int is not a fixed-size type
    var x int
    binary.Write(&buf, binary.LittleEndian, x)
}
```

### Valid

```golang
package main

import (
    "bytes"
    "encoding/binary"
)

func main() {
    var buf bytes.Buffer
    // int32 is a fixed-size type
    var x int32
    binary.Write(&buf, binary.LittleEndian, x)
}
```
