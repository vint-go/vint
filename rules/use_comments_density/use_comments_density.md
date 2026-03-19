---
title: useCommentsDensity
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useCommentsDensity`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useCommentsDensity:
    arguments:
      - 15
```

## Details

Spots files not respecting a minimum value for the comments lines density metric. The density is calculated as: comment lines / (lines of code + comment lines) * 100.

The rule accepts an integer argument specifying the minimum expected comments lines density percentage. If no argument is provided, the default minimum is 0 (effectively disabled).

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With a minimum density of 60%, this file would fail because
// the comment density is only about 57%.
package fixtures

// func does something
func cd1() error {
    // the var
    var x string
    /* the return */
    return nil
}
```

### Valid

```golang
// With a minimum density of 15%, this file would pass because
// the comment density exceeds the threshold.
package fixtures

// wellDocumented is a function that demonstrates good comment density.
// It has sufficient documentation to meet the minimum threshold.
func wellDocumented() error {
    // initialize
    var x string
    // return result
    return nil
}
```
