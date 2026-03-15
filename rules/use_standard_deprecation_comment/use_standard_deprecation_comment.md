---
title: useStandardDeprecationComment
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useStandardDeprecationComment`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useStandardDeprecationComment:
    # no additional options
```

## Details

Detects malformed 'deprecated' doc-comments. The Go standard requires deprecation notices to follow the format `// Deprecated: <explanation>`. This checker identifies issues including:

- Incorrect casing (e.g., "DEPRECATED:" instead of "Deprecated:")
- Using comma instead of colon ("Deprecated," vs "Deprecated:")
- Alternative patterns like "this function is deprecated" or "deprecated. use"
- Typographical errors in the prefix word
- Deprecation notices appearing inline rather than in dedicated paragraphs

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
// DEPRECATED: use NewFunc instead
func OldFunc() {}
```

```golang
// this function is deprecated, use NewFunc
func OldFunc() {}
```

### Valid

```golang
// Deprecated: use NewFunc instead.
func OldFunc() {}
```
