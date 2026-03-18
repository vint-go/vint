---
title: noUnnecessaryBlankIdentifier
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnnecessaryBlankIdentifier`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnnecessaryBlankIdentifier:
    # rule options here
```

## Details

Drop unnecessary use of the blank identifier.

When only the value from a range is needed, the blank identifier for the index is unnecessary and can be omitted.

Source: https://staticcheck.dev/docs/checks/#S1005

## Examples

### Invalid

```golang
package main

import "fmt"

func process(items []string) {
    for _ = range items {
        fmt.Println("item")
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(items []string) {
    for range items {
        fmt.Println("item")
    }
}
```
