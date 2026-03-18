---
title: useRawStringRegexp
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useRawStringRegexp`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useRawStringRegexp:
    # rule options here
```

## Details

Simplify regular expression by using raw string literal.

Regular expressions containing backslash escapes are easier to read as raw string literals (backtick-quoted) because backslashes don't need to be doubled.

Source: https://staticcheck.dev/docs/checks/#S1007

## Examples

### Invalid

```golang
package main

import "regexp"

func main() {
    // Double backslashes needed in interpreted string
    re := regexp.MustCompile("\\d+\\.\\d+")
    _ = re
}
```

### Valid

```golang
package main

import "regexp"

func main() {
    // Raw string literal - no double backslashes needed
    re := regexp.MustCompile(`\d+\.\d+`)
    _ = re
}
```
