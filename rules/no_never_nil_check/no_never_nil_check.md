---
title: noNeverNilCheck
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noNeverNilCheck`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noNeverNilCheck:
    # rule options here
```

## Details

Checking never-nil value against nil.

Checking a value that can never be nil against nil is pointless. This typically occurs with composite literals, return values of functions known to never return nil, or other expressions whose type guarantees a non-nil value.

Source: https://staticcheck.dev/docs/checks/#SA4031

## Examples

### Invalid

```golang
package main

func process() {
    m := map[string]int{}
    // m is a composite literal, never nil
    if m == nil {
        panic("impossible")
    }
}
```

### Valid

```golang
package main

func process(m map[string]int) {
    // m is a parameter, could be nil
    if m == nil {
        m = make(map[string]int)
    }
}
```
