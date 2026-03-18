---
title: useIdiomaticDurationName
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useIdiomaticDurationName`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useIdiomaticDurationName:
    # rule options here
```

## Details

Poorly chosen name for variable of type `time.Duration`.

Variables of type `time.Duration` should not have a name that implies a specific time unit, such as `timeoutSecs` or `delayMs`. Since `time.Duration` already represents a duration with its own unit, names like `timeout` or `delay` are more appropriate.

Source: https://staticcheck.dev/docs/checks/#ST1011

## Examples

### Invalid

```golang
package main

import "time"

// Name implies seconds, but Duration has its own unit
var timeoutSecs time.Duration = 5 * time.Second
```

### Valid

```golang
package main

import "time"

var timeout time.Duration = 5 * time.Second
```
