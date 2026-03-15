---
title: noCommentedOutCode
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noCommentedOutCode`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noCommentedOutCode:
    # no additional options
```

## Details

Detects commented-out code inside function bodies. This checker uses heuristics to identify comments that appear to contain executable Go code rather than actual documentation. Comments should describe code, not contain disabled code that should be removed or managed via version control.

The checker applies several filters to minimize false positives:
- Skips comments containing markers like "TODO", URLs, or explanatory phrases like "e.g."
- Ignores very short comments (default minimum length: 15 characters)
- Allows certain statement types and patterns

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
func process() {
    // fmt.Println("debug output")
    // result := computeValue()
    doWork()
}
```

### Valid

```golang
func process() {
    // Process the work items and return results
    doWork()
}
```
