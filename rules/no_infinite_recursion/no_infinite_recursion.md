---
title: noInfiniteRecursion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noInfiniteRecursion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noInfiniteRecursion:
    # rule options here
```

## Details

Infinite recursive call.

A function that calls itself unconditionally in all code paths will cause a stack overflow. This detects direct and indirect infinite recursion.

Source: https://staticcheck.dev/docs/checks/#SA5007

## Examples

### Invalid

```golang
package main

func factorial(n int) int {
    // Missing base case - infinite recursion
    return n * factorial(n-1)
}
```

### Valid

```golang
package main

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}
```
