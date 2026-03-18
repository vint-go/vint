---
title: useTaggedSwitch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useTaggedSwitch`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a **fix** available.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useTaggedSwitch:
    # rule options here
```

## Details

Convert untagged switch to tagged switch.

An untagged switch with comparisons against the same variable can be rewritten as a tagged switch, which is often clearer.

Source: https://staticcheck.dev/docs/checks/#QF1002

## Examples

### Invalid

```golang
package main

func process(x int) string {
    switch {
    case x == 1:
        return "one"
    case x == 2:
        return "two"
    default:
        return "other"
    }
}
```

### Valid

```golang
package main

func process(x int) string {
    switch x {
    case 1:
        return "one"
    case 2:
        return "two"
    default:
        return "other"
    }
}
```
