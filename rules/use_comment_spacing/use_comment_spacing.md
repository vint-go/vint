---
title: useCommentSpacing
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCommentSpacing`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCommentSpacing:
    # no additional options
```

## Details

Detects comments with non-idiomatic formatting. Go convention requires single-line comments to have a space after the `//` prefix. This checker enforces that formatting.

The checker exempts special cases including:
- Block comments (`/* */`)
- Compiler directives (`//go:generate`, `//line`)
- Linting directives (`//nolint`, `//noinspection`)
- IDE folding regions (`//region`, `//editor-fold`)
- Custom markers starting with `+`, `-`, `#`, or `!`
- Keyed directives matching `//key:` patterns (e.g., `//TODO:`, `//FIXME:`, `//nolint:foo`)

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
//This comment has no space after the slashes
func f() {}
```

### Valid

```golang
// This comment has proper spacing
func f() {}
```

```golang
// TODO: fix this later
```

```golang
//go:generate stringer -type=Status
```

```golang
//TODO: fix this later
```
