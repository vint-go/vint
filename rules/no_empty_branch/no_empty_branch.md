---
title: noEmptyBranch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noEmptyBranch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noEmptyBranch:
    # no additional options
```

## Details

Empty body in an `if` or `else` branch.

Having an empty body in an `if` or `else` clause usually indicates incomplete code. If intentional, add a comment explaining why.

Source: https://staticcheck.dev/docs/checks/#SA9003

## Examples

### Invalid

```golang
package main

func process(x int) {
    if x > 0 {
        // Empty body
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(x int) {
    if x > 0 {
        fmt.Println("positive")
    }
}
```
