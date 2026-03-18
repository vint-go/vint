---
title: noInvalidRegexp
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidRegexp`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidRegexp:
    # rule options here
```

## Details

Invalid regular expression.

This check validates that arguments to `regexp.Compile`, `regexp.MustCompile`, `regexp.Match`, and related functions are valid regular expressions. Providing an invalid regular expression will cause a runtime panic (for `MustCompile`) or an error at runtime.

Source: https://staticcheck.dev/docs/checks/#SA1000

## Examples

### Invalid

```golang
package main

import "regexp"

func main() {
    // Invalid regex - unclosed group
    re := regexp.MustCompile("foo(bar")
    _ = re
}
```

### Valid

```golang
package main

import "regexp"

func main() {
    // Valid regex
    re := regexp.MustCompile("foo(bar)")
    _ = re
}
```
