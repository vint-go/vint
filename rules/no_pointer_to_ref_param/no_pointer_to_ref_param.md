---
title: noPointerToRefParam
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noPointerToRefParam`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noPointerToRefParam:
    # no additional options
```

## Details

Detects input and output parameters that have a type of pointer to referential type. Maps, channels, and interfaces already have reference semantics, so wrapping them in a pointer is redundant. This checker flags both named and unnamed parameters where unnecessary pointer indirection is used.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func f(m *map[string]int) (*chan *int) {
    // pointer to map and pointer to channel are redundant
}
```

### Valid

```golang
func f(m map[string]int) (chan *int) {
    // map and channel used without unnecessary pointer
}
```
