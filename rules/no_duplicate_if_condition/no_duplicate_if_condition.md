---
title: noDuplicateIfCondition
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDuplicateIfCondition`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDuplicateIfCondition:
    # rule options here
```

## Details

An `if/else if` chain has repeated conditions and therefore always takes the same branch.

When an `if/else if` chain contains duplicate conditions, the second occurrence can never be reached because the first condition already matched. This is usually a copy-paste error.

Source: https://staticcheck.dev/docs/checks/#SA4014

## Examples

### Invalid

```golang
package main

func check(x int) string {
    if x > 10 {
        return "big"
    } else if x > 10 {
        // Unreachable: same condition as above
        return "also big"
    }
    return "small"
}
```

### Valid

```golang
package main

func check(x int) string {
    if x > 10 {
        return "big"
    } else if x > 5 {
        return "medium"
    }
    return "small"
}
```
