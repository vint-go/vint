---
title: noDuplicateCutsetChars
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateCutsetChars`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateCutsetChars:
    # rule options here
```

## Details

A string `cutset` passed to `strings.TrimLeft` or `strings.TrimRight` contains duplicate characters.

`strings.TrimLeft` and `strings.TrimRight` treat their second argument as a set of characters, not a substring. Duplicate characters in the cutset are redundant and may indicate that `strings.TrimPrefix` or `strings.TrimSuffix` was intended instead.

Source: https://staticcheck.dev/docs/checks/#SA1024

## Examples

### Invalid

```golang
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Likely intended TrimPrefix, not TrimLeft
    s := strings.TrimLeft("httpexample.com", "http")
    fmt.Println(s) // Removes all 'h', 't', 'p' characters from the left
}
```

### Valid

```golang
package main

import (
    "fmt"
    "strings"
)

func main() {
    // Use TrimPrefix to remove a prefix string
    s := strings.TrimPrefix("http://example.com", "http://")
    fmt.Println(s) // "example.com"
}
```
