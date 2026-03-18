---
title: noNilMapAssignment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noNilMapAssignment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noNilMapAssignment:
    # rule options here
```

## Details

Assignment to nil map.

Assigning a value to a nil map causes a runtime panic. Maps must be initialized with `make` or a composite literal before use.

Source: https://staticcheck.dev/docs/checks/#SA5000

## Examples

### Invalid

```golang
package main

func main() {
    var m map[string]int
    // Panic: assignment to entry in nil map
    m["key"] = 1
}
```

### Valid

```golang
package main

func main() {
    m := make(map[string]int)
    m["key"] = 1
}
```
