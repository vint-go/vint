---
title: noOddSizeSliceArg
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noOddSizeSliceArg`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noOddSizeSliceArg:
    # rule options here
```

## Details

Passing odd-sized slice to function expecting even size.

Some functions expect slices with an even number of elements (e.g., key-value pairs). Passing a slice with an odd number of elements indicates a logic error.

Source: https://staticcheck.dev/docs/checks/#SA5012

## Examples

### Invalid

```golang
package main

import "log/slog"

func main() {
    // Odd number of arguments - missing a value
    slog.Info("message", "key1", "value1", "key2")
}
```

### Valid

```golang
package main

import "log/slog"

func main() {
    // Even number of arguments - key-value pairs
    slog.Info("message", "key1", "value1", "key2", "value2")
}
```
