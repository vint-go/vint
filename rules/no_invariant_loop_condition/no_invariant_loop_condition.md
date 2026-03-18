---
title: noInvariantLoopCondition
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noInvariantLoopCondition`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noInvariantLoopCondition:
    # rule options here
```

## Details

The variable in the loop condition never changes, leading to an infinite or zero-iteration loop.

If the condition variable in a `for` loop is never modified within the loop body, the loop will either never execute or never terminate.

Source: https://staticcheck.dev/docs/checks/#SA4008

## Examples

### Invalid

```golang
package main

func process(items []string) {
    i := 0
    for i < len(items) {
        // i is never incremented, infinite loop
        process(items)
    }
}
```

### Valid

```golang
package main

import "fmt"

func process(items []string) {
    for i := 0; i < len(items); i++ {
        fmt.Println(items[i])
    }
}
```
