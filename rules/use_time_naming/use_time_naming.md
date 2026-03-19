---
title: useTimeNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTimeNaming`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTimeNaming:
    # no configuration options
```

## Details

Using unit-specific suffix like "Secs", "Mins", ... when naming variables of type `time.Duration` can be misleading, this rule highlights those cases.

Since `time.Duration` already represents a duration with its own unit, adding suffixes that imply a specific time unit (such as "Hour", "Sec", "Ms", etc.) is redundant and can lead to confusion. Use names like `timeout` or `delay` instead of `timeoutSecs` or `delayMs`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package main

import "time"

var timeoutSecs = 5 * time.Second    // don't use unit-specific suffix "Secs"
var delayMs = 100 * time.Millisecond // don't use unit-specific suffix "Ms"
var retryMin = 2 * time.Minute      // don't use unit-specific suffix "Min"
```

### Valid

```golang
package main

import "time"

var timeout = 5 * time.Second
var delay = 100 * time.Millisecond
var retryInterval = 2 * time.Minute
```
