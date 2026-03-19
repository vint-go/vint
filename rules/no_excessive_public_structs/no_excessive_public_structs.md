---
title: noExcessivePublicStructs
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noExcessivePublicStructs`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noExcessivePublicStructs:
    arguments: [3]
```

## Details

Packages declaring too many public structs can be hard to understand/use,
and could be a symptom of bad design.

This rule warns on files declaring more than a configured maximum number of public struct declarations.
The default maximum is 5.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With max set to 1, having more than 1 public struct triggers the rule.
package pkg

type Foo struct {}

type Bar struct {} // exceeds the limit of 1 public struct
```

### Valid

```golang
// With default max of 5, having 3 public structs is fine.
package pkg

type Foo struct {}

type Bar struct {}

type Baz struct {}
```
