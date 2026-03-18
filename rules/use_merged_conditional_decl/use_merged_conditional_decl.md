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

Merge conditional assignment into variable declaration.

A variable declaration followed by a conditional assignment can be merged into a single declaration with a conditional expression.

Source: https://staticcheck.dev/docs/checks/#QF1007

## Examples

### Invalid

```golang
package main

func process(useDefault bool) string {
    x := "custom"
    if useDefault {
        x = "default"
    }
    return x
}
```

### Valid

```golang
package main

func process(useDefault bool) string {
    if useDefault {
        return "default"
    }
    return "custom"
}
```
