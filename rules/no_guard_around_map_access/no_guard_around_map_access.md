---
title: noGuardAroundMapAccess
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noGuardAroundMapAccess`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noGuardAroundMapAccess:
    # rule options here
```

## Details

Unnecessary guard around map access.

Checking `if _, ok := m[k]; ok { v = m[k] }` can be simplified to `v = m[k]` since accessing a missing key returns the zero value.

Source: https://staticcheck.dev/docs/checks/#S1036

## Examples

### Invalid

```golang
package main

func getValue(m map[string]int, key string) int {
    var v int
    if _, ok := m[key]; ok {
        v = m[key]
    }
    return v
}
```

### Valid

```golang
package main

func getValue(m map[string]int, key string) int {
    return m[key]
}
```
