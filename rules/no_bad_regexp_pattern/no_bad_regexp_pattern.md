---
title: noBadRegexpPattern
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noBadRegexpPattern`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noBadRegexpPattern:
    # no additional options
```

## Details

Detects suspicious regexp patterns that may contain errors. This checker analyzes regular expressions passed to `regexp.Compile` and `regexp.MustCompile` for several classes of issues:

1. **Dangling/Redundant Anchors** -- Detects `^` or `$` positioned improperly or used only on specific alternatives.
2. **Nested Quantifiers** -- Identifies repeated greedy quantifiers like `a**` or `a*+`.
3. **Duplicated Alternations** -- Finds identical options in alternation groups (e.g., `a|b|a`).
4. **Suspicious Character Ranges** -- Flags non-standard ranges like `!-_` that aren't numeric or letter sequences.
5. **Character Class Overlaps** -- Detects duplicate or intersecting character definitions within classes.
6. **Redundant/Contradictory Flags** -- Reports setting flags twice or clearing unset flags.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// Anchor inside non-capturing group
re := regexp.MustCompile(`(?:^foo|bar)`)
```

```golang
// Overlapping character class
re := regexp.MustCompile(`[aba]`)
```

### Valid

```golang
// Anchor outside group
re := regexp.MustCompile(`^(?:foo|bar)`)
```

```golang
// No overlapping characters
re := regexp.MustCompile(`[ab]`)
```
