---
title: useStringsContains
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useStringsContains`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useStringsContains:
    # rule options here
```

## Details

Replace call to `strings.Index` with `strings.Contains`.

Using `strings.Index(s, substr) != -1` or `>= 0` to check for substring presence can be simplified to `strings.Contains(s, substr)`.

Source: https://staticcheck.dev/docs/checks/#S1003

## Examples

### Invalid

```golang
package main

import "strings"

func hasPrefix(s string) bool {
    return strings.Index(s, "hello") != -1
}
```

### Valid

```golang
package main

import "strings"

func hasPrefix(s string) bool {
    return strings.Contains(s, "hello")
}
```
