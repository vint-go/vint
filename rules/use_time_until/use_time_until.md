---
title: useTimeUntil
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTimeUntil`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTimeUntil:
    # rule options here
```

## Details

Replace `x.Sub(time.Now())` with `time.Until(x)`.

`time.Until(x)` is a more idiomatic and readable way to express `x.Sub(time.Now())`.

Source: https://staticcheck.dev/docs/checks/#S1024

## Examples

### Invalid

```golang
package main

import "time"

func remaining(deadline time.Time) time.Duration {
    return deadline.Sub(time.Now())
}
```

### Valid

```golang
package main

import "time"

func remaining(deadline time.Time) time.Duration {
    return time.Until(deadline)
}
```
