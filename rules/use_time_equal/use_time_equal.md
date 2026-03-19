---
title: useTimeEqual
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/useTimeEqual`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/useTimeEqual:
    # no configuration options
```

## Details

Suggests using `time.Time.Equal` instead of `==` and `!=` operators for equality checks on `time.Time` values.

`time.Time` values should be compared using the `Equal` method rather than `==`, because `==` also compares the monotonic clock reading and the location, which may differ even for the same instant in time. For more information see the [time.Time documentation](https://pkg.go.dev/time#Time).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package main

import "time"

func isSame(a, b time.Time) bool {
    return a == b // use a.Equal(b) instead
}

func isDifferent(a, b time.Time) bool {
    return a != b // use !a.Equal(b) instead
}
```

### Valid

```golang
package main

import "time"

func isSame(a, b time.Time) bool {
    return a.Equal(b)
}

func isDifferent(a, b time.Time) bool {
    return !a.Equal(b)
}
```
