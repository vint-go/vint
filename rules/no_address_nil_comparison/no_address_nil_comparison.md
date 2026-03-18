---
title: noAddressNilComparison
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noAddressNilComparison`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noAddressNilComparison:
    # rule options here
```

## Details

Comparing the address of a variable against nil.

The address of a variable is never nil. Comparisons like `&x == nil` are always false, and `&x != nil` are always true. This is likely a logic error.

Source: https://staticcheck.dev/docs/checks/#SA4022

## Examples

### Invalid

```golang
package main

func process() {
    var x int
    // Address of a variable is never nil
    if &x == nil {
        // unreachable
    }
}
```

### Valid

```golang
package main

func process(x *int) {
    // Checking a pointer parameter is meaningful
    if x == nil {
        return
    }
}
```
