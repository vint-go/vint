---
title: noUnescapedRegexpDot
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnescapedRegexpDot`
- This rule is not enabled by default and must be explicitly enabled.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnescapedRegexpDot:
    # no additional options
```

## Details

Detects suspicious regexp patterns with unescaped dots before common domain extensions. This checker identifies unescaped dots in regular expression patterns used with Go's `regexp` package functions. Specifically, it detects when common domain extensions (com, org, info, net, etc.) appear after a non-escaped period, since the dot matches any character instead of a literal period, which is likely unintended in domain-matching contexts.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
re := regexp.MustCompile(`google.com`) // dot matches any character
```

### Valid

```golang
re := regexp.MustCompile(`google\.com`) // dot is properly escaped
```
