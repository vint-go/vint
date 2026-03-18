---
title: noInvalidUtf8StringArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidUtf8StringArg`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidUtf8StringArg:
    # rule options here
```

## Details

Various methods in the `strings` package expect valid UTF-8, but invalid input is provided.

Functions like `strings.ToUpper`, `strings.ToLower`, and similar expect valid UTF-8 encoded strings. Passing invalid UTF-8 data may lead to unexpected results.

Source: https://staticcheck.dev/docs/checks/#SA1011

## Examples

### Invalid

```golang
package main

import "strings"

func main() {
    // Invalid UTF-8 byte sequence
    s := string([]byte{0xff, 0xfe})
    result := strings.ToUpper(s)
    _ = result
}
```

### Valid

```golang
package main

import "strings"

func main() {
    // Valid UTF-8 string
    s := "hello world"
    result := strings.ToUpper(s)
    _ = result
}
```
