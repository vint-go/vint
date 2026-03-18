---
title: useTimeSince
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTimeSince`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTimeSince:
    # rule options here
```

## Details

Replace `time.Now().Sub(x)` with `time.Since(x)`.

`time.Since(x)` is a more idiomatic and readable way to express `time.Now().Sub(x)`.

Source: https://staticcheck.dev/docs/checks/#S1012

## Examples

### Invalid

```golang
package main

import "time"

func elapsed(start time.Time) time.Duration {
    return time.Now().Sub(start)
}
```

### Valid

```golang
package main

import "time"

func elapsed(start time.Time) time.Duration {
    return time.Since(start)
}
```
