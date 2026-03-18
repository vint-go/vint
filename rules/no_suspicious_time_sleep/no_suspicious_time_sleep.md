---
title: noSuspiciousTimeSleep
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noSuspiciousTimeSleep`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noSuspiciousTimeSleep:
    # rule options here
```

## Details

Suspiciously small untyped constant in `time.Sleep`.

`time.Sleep` expects a `time.Duration` value. Passing a small integer literal like `1` means sleeping for 1 nanosecond, which is almost certainly not the intended behavior. You likely meant to multiply by a time unit such as `time.Second`.

Source: https://staticcheck.dev/docs/checks/#SA1004

## Examples

### Invalid

```golang
package main

import "time"

func main() {
    // Sleeps for 1 nanosecond, almost certainly not intended
    time.Sleep(1)
}
```

### Valid

```golang
package main

import "time"

func main() {
    // Sleeps for 1 second
    time.Sleep(1 * time.Second)
}
```
