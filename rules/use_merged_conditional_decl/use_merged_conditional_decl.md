---
title: useMergedConditionalDecl
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useMergedConditionalDecl`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useMergedConditionalDecl:
    # rule options here
```

## Details

Merge conditional boolean assignment into variable declaration.

A boolean variable declaration initialized to `true` or `false`, followed by a
conditional reassignment to the opposite boolean literal, can be simplified into
a single declaration using the condition directly.

This rule only applies to **boolean literals**. Non-boolean types (strings,
integers, slices, etc.) are not flagged.

Source: https://staticcheck.dev/docs/checks/#QF1007

## Examples

### Invalid

```golang
package main

func process(cond bool) bool {
    x := false
    if cond {
        x = true
    }
    return x
}
// Can be simplified to: x := cond
```

```golang
package main

func process(cond bool) bool {
    x := true
    if cond {
        x = false
    }
    return x
}
// Can be simplified to: x := !cond
```

### Valid

```golang
package main

// Non-boolean types are not flagged.
func process(useDefault bool) string {
    x := "custom"
    if useDefault {
        x = "default"
    }
    return x
}
```
