---
title: useTimeSleep
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTimeSleep`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTimeSleep:
    # rule options here
```

## Details

Elaborate way of sleeping.

Using `select { case <-time.After(d): }` is a more complex way of writing `time.Sleep(d)`.

Source: https://staticcheck.dev/docs/checks/#S1037

## Examples

### Invalid

```golang
package main

import "time"

func wait() {
    select {
    case <-time.After(1 * time.Second):
    }
}
```

### Valid

```golang
package main

import "time"

func wait() {
    time.Sleep(1 * time.Second)
}
```
