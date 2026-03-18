---
title: useTimeEqual
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useTimeEqual`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useTimeEqual:
    # rule options here
```

## Details

Use `time.Time.Equal` instead of `==` operator.

`time.Time` values should be compared using the `Equal` method rather than `==`, because `==` also compares the location, which may differ even for the same instant in time.

Source: https://staticcheck.dev/docs/checks/#QF1009

## Examples

### Invalid

```golang
package main

import "time"

func isSame(a, b time.Time) bool {
    return a == b
}
```

### Valid

```golang
package main

import "time"

func isSame(a, b time.Time) bool {
    return a.Equal(b)
}
```
