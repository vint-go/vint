---
title: noRedundantNilLoopCheck
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantNilLoopCheck`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantNilLoopCheck:
    # rule options here
```

## Details

Omit redundant nil check around loop.

Checking `if s != nil` before `for range s` is redundant because ranging over a nil slice or map simply does nothing (zero iterations).

Source: https://staticcheck.dev/docs/checks/#S1031

## Examples

### Invalid

```golang
package main

import "fmt"

func process(items []string) {
    if items != nil {
        for _, item := range items {
            fmt.Println(item)
        }
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(items []string) {
    for _, item := range items {
        fmt.Println(item)
    }
}
```
