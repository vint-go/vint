---
title: noIneffectiveRandCall
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noIneffectiveRandCall`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noIneffectiveRandCall:
    # rule options here
```

## Details

Ineffective attempt at generating random number.

Using `rand.Intn(1)` always returns 0 because the range `[0, 1)` contains only the value 0. This is likely a logic error.

Source: https://staticcheck.dev/docs/checks/#SA4030

## Examples

### Invalid

```golang
package main

import "math/rand"

func getRandomIndex() int {
    // Always returns 0
    return rand.Intn(1)
}
```

### Valid

```golang
package main

import "math/rand"

func getRandomIndex(max int) int {
    return rand.Intn(max)
}
```
