---
title: noMalformedDirective
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noMalformedDirective`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noMalformedDirective:
    # rule options here
```

## Details

Checks known Go toolchain directives. This analyzer validates that Go compiler directives (comments starting with `//go:`) are correctly formed and placed. It checks for:

- Misspelled directives (e.g., `//go:noinlin` instead of `//go:noinline`)
- Directives placed in the wrong position (e.g., `//go:generate` not at the beginning of a line)
- Directives with invalid arguments
- Unknown directives that look like they might be typos of known ones

Source: https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/directive

## Examples

### Invalid

```golang
// go:noinline
// Bad: space between // and go: makes it a regular comment, not a directive
func example() {}
```

```golang
//go:noinlin
// Bad: misspelled directive (missing 'e')
func example() {}
```

### Valid

```golang
//go:noinline
func example() {}
```

```golang
//go:generate stringer -type=MyType
```
