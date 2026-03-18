---
title: useDirectStringRange
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useDirectStringRange`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useDirectStringRange:
    # rule options here
```

## Details

Range over string directly instead of converting to rune slice.

Ranging over a string already yields runes. Converting to `[]rune` before ranging is unnecessary.

Source: https://staticcheck.dev/docs/checks/#S1029

## Examples

### Invalid

```golang
package main

import "fmt"

func process(s string) {
    for _, r := range []rune(s) {
        fmt.Printf("%c", r)
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(s string) {
    for _, r := range s {
        fmt.Printf("%c", r)
    }
}
```
