---
title: noRedundantRuneConversion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/noRedundantRuneConversion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/noRedundantRuneConversion:
    # rule options here
```

## Details

Converting a string to a slice of runes before ranging over it is unnecessary.

`for _, r := range string` already iterates over runes. Converting to `[]rune` first is unnecessary and incurs an allocation.

Source: https://staticcheck.dev/docs/checks/#SA6003

## Examples

### Invalid

```golang
package main

import "fmt"

func process(s string) {
    // Unnecessary conversion to []rune
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
    // Ranging over string already iterates runes
    for _, r := range s {
        fmt.Printf("%c", r)
    }
}
```
