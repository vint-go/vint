---
title: useCommentSpacings
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCommentSpacings`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCommentSpacings:
    arguments:
      - "mypragma:"
      - "+optional"
```

## Details

Spots comments that lack a space between the comment delimiter (`//` or `/*`) and the start of the comment text. This rule is configurable: you can provide a list of string prefixes to exempt from the check (e.g. custom directives like `mypragma:` or markers like `+optional`). Standard Go directives such as `//go:generate`, `//nolint`, `//export`, `//extern`, and `//line` are automatically exempted.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
//This comment does not respect the spacing rule!
var a string
```

```golang
/*Not valid
 */
```

### Valid

```golang
// This is a well formed comment
var a string
```

```golang
//nolint:staticcheck // nolint is in the default list of acceptable comments.
var b string
```

```golang
//go:generate stringer -type=Pill
```

```golang
/*
Should be valid
*/
```
