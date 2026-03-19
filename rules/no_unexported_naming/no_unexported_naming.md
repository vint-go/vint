---
title: noUnexportedNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noUnexportedNaming`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noUnexportedNaming:
    # no configuration options
```

## Details

Warns on wrongly named unexported symbols, i.e. unexported symbols whose name starts with a capital letter. In Go, local variables, function parameters, and return values should start with a lowercase letter since they cannot be exported. Using a capitalized name for such symbols is misleading and violates Go naming conventions.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
func example(
    S int, // the symbol S is local, its name should start with a lowercase letter
    s int,
) (
    Result bool, // the symbol Result is local, its name should start with a lowercase letter
    result bool,
) {
    var NotExportable int // the symbol NotExportable is local, its name should start with a lowercase letter
    return
}
```

### Valid

```golang
func example(
    s int,
    t int,
) (
    result bool,
    ok bool,
) {
    var local int
    return
}
```
