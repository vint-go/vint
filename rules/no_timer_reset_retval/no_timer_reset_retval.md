---
title: noTimerResetRetval
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noTimerResetRetval`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noTimerResetRetval:
    # rule options here
```

## Details

It is not possible to use `(*time.Timer).Reset`'s return value correctly.

The return value of `(*time.Timer).Reset` is unreliable and should not be used to determine whether the timer had already fired. The documentation explicitly warns against using this return value for synchronization.

Source: https://staticcheck.dev/docs/checks/#SA1025

## Examples

### Invalid

```golang
package main

import "time"

func main() {
    t := time.NewTimer(1 * time.Second)
    // Wrong: using Reset's return value
    if !t.Reset(2 * time.Second) {
        <-t.C
    }
}
```

### Valid

```golang
package main

import "time"

func main() {
    t := time.NewTimer(1 * time.Second)
    if !t.Stop() {
        <-t.C
    }
    t.Reset(2 * time.Second)
}
```
