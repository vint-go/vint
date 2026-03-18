---
title: noRedundantLabel
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantLabel`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantLabel:
    # no additional options
```

## Details

Detects redundant statement labels. This checker identifies two patterns of redundant labels in Go code:

1. **Completely redundant labels** -- Labels on statements where no nested break/continue statements reference that label.
2. **Labeled continue anti-pattern** -- Cases where a labeled `continue` statement targeting an outer loop could be simplified.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
loop:
for i := 0; i < n; i++ {
    // no break or continue references "loop"
    doSomething(i)
}
```

### Valid

```golang
for i := 0; i < n; i++ {
    doSomething(i)
}
```

```golang
outer:
for i := 0; i < n; i++ {
    for j := 0; j < m; j++ {
        if condition {
            break outer // label is actually used
        }
    }
}
```
