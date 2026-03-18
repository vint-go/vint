---
title: useSimplifiedPrintFormat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSimplifiedPrintFormat`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSimplifiedPrintFormat:
    # rule options here
```

## Details

Unnecessarily complex way of printing formatted string.

Using `fmt.Print(fmt.Sprintf(...))` can be simplified to `fmt.Printf(...)`.

Source: https://staticcheck.dev/docs/checks/#S1038

## Examples

### Invalid

```golang
package main

import "fmt"

func main() {
    fmt.Println(fmt.Sprintf("Hello, %s!", "world"))
}
```

### Valid

```golang
package main

import "fmt"

func main() {
    fmt.Printf("Hello, %s!\n", "world")
}
```
