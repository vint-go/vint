---
title: noRedundantControlFlow
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantControlFlow`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantControlFlow:
    # rule options here
```

## Details

Omit redundant control flow.

A `return` at the end of a function that returns nothing, or a `break` at the end of a case clause, is redundant and can be removed.

Source: https://staticcheck.dev/docs/checks/#S1023

## Examples

### Invalid

```golang
package main

import "fmt"

func process() {
    fmt.Println("done")
    return // redundant
}
```

### Valid

```golang
package main

import "fmt"

func process() {
    fmt.Println("done")
}
```
