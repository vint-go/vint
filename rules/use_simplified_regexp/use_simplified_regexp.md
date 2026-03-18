---
title: useSimplifiedRegexp
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useSimplifiedRegexp`
- This rule is not recommended (experimental, opinionated).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useSimplifiedRegexp:
    # no additional options
```

## Details

Detects regexp patterns that can be simplified. This checker analyzes regular expression patterns for simplification opportunities, including:

- Converting verbose quantifiers like `{0,1}`, `{1,}`, `{0,}` to `?`, `+`, `*`
- Replacing character classes (e.g., `[0-9]` to `\d`, `[[:space:]]` to `\s`)
- Simplifying negated character classes (e.g., `[^0-9]` to `\D`)
- Combining repeated single characters into quantified forms
- Converting alternations of single characters to character classes (`a|b|c` to `[abc]`)
- Factoring common prefixes/suffixes in alternations (`http|https` to `https?`)
- Removing unnecessary non-capturing group wrappers

The checker skips patterns exceeding 60 characters and performs up to 2 simplification passes.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
re := regexp.MustCompile(`[0-9]+`)
```

```golang
re := regexp.MustCompile(`http|https`)
```

### Valid

```golang
re := regexp.MustCompile(`\d+`)
```

```golang
re := regexp.MustCompile(`https?`)
```
