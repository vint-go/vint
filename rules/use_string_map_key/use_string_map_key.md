---
title: useStringMapKey
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/performance/useStringMapKey`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/performance/useStringMapKey:
    # rule options here
```

## Details

Missing an optimization opportunity when indexing maps by byte slices.

When using a `map[string]` and you have a `[]byte`, converting to `string` for map lookup allocates a new string. The compiler can optimize `m[string(b)]` but not intermediate variables. Use `m[string(b)]` directly.

Source: https://staticcheck.dev/docs/checks/#SA6001

## Examples

### Invalid

```golang
package main

func lookup(m map[string]int, key []byte) int {
    // Intermediate variable prevents optimization
    s := string(key)
    return m[s]
}
```

### Valid

```golang
package main

func lookup(m map[string]int, key []byte) int {
    // Direct conversion allows compiler optimization
    return m[string(key)]
}
```
