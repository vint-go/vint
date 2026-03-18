---
title: noDubiousBitShift
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noDubiousBitShift`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noDubiousBitShift:
    # rule options here
```

## Details

Dubious bit shifting of a fixed size integer value.

Shifting an integer by more bits than its size always results in zero (for unsigned) or zero/-1 (for signed). For example, shifting a `uint8` by 8 or more bits always yields zero.

Source: https://staticcheck.dev/docs/checks/#SA9006

## Examples

### Invalid

```golang
package main

func process() uint8 {
    var x uint8 = 1
    // Shifting uint8 by 8 bits always yields 0
    return x << 8
}
```

### Valid

```golang
package main

func process() uint16 {
    var x uint16 = 1
    // uint16 can hold the result of shifting by 8
    return x << 8
}
```
