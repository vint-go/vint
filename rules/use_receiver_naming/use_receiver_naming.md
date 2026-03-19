---
title: useReceiverNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useReceiverNaming`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useReceiverNaming:
    arguments:
      - maxLength: 2
```

## Details

By convention, receiver names in a method should reflect their identity.
For example, if the receiver is of type `Parts`, `p` is an adequate name for it.
Contrary to other languages, it is not idiomatic to name receivers as `this` or `self`.

The rule also checks that:
- Receiver names are not underscores (omit the name instead if unused).
- Receiver names are consistent across all methods of the same type.
- Receiver names do not exceed a configurable maximum length.

Configuration options:
- `max-length` (int): maximum allowed length for receiver names. When set, receiver names longer than this value will be flagged.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// Using generic names like "this" or "self"
func (this *MyStruct) Method() {}
func (self *MyStruct) Method() {}

// Using underscore as receiver name
func (_ *MyStruct) Method() {}

// Inconsistent receiver names for the same type
func (m *MyStruct) Method1() {}
func (s *MyStruct) Method2() {} // should use "m" consistently
```

### Valid

```golang
// Short, idiomatic receiver name reflecting identity
func (m *MyStruct) Method1() {}
func (m *MyStruct) Method2() {}

// Consistent receiver naming
func (p *Parts) Assemble() {}
func (p *Parts) Disassemble() {}
```
