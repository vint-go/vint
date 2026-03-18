---
title: noInvalidStrconvArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInvalidStrconvArg`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInvalidStrconvArg:
    # rule options here
```

## Details

Invalid argument in call to `strconv` function.

This check detects invalid arguments passed to `strconv.ParseInt`, `strconv.ParseUint`, `strconv.ParseFloat`, and `strconv.FormatInt`. For example, an invalid base or bit size argument.

Source: https://staticcheck.dev/docs/checks/#SA1030

## Examples

### Invalid

```golang
package main

import "strconv"

func main() {
    // Invalid base: base must be 0 or between 2 and 36
    _, _ = strconv.ParseInt("42", 1, 64)
}
```

### Valid

```golang
package main

import "strconv"

func main() {
    // Valid base 10
    _, _ = strconv.ParseInt("42", 10, 64)
}
```
