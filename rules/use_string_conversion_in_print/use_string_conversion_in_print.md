---
title: useStringConversionInPrint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useStringConversionInPrint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useStringConversionInPrint:
    # rule options here
```

## Details

Convert slice of bytes to string when printing with `fmt`.

Passing `[]byte` to `fmt.Println` and similar functions prints the byte values. If you want the string representation, convert to string first.

Source: https://staticcheck.dev/docs/checks/#QF1010

## Examples

### Invalid

```golang
package main

import "fmt"

func main() {
    b := []byte("hello")
    // Prints byte values: [104 101 108 108 111]
    fmt.Println(b)
}
```

### Valid

```golang
package main

import "fmt"

func main() {
    b := []byte("hello")
    // Prints the string: hello
    fmt.Println(string(b))
}
```
