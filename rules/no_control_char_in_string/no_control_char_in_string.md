---
title: noControlCharInString
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noControlCharInString`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noControlCharInString:
    # no additional options
```

## Details

Avoid zero-width and control characters in string literals.

String literals should not contain invisible Unicode characters such as zero-width spaces, zero-width joiners, or other control characters. These characters can be confusing and may cause subtle bugs.

Source: https://staticcheck.dev/docs/checks/#ST1018

## Examples

### Invalid

```golang
// String contains invisible zero-width space character
var s = "hello\u200bworld"
```

### Valid

```golang
var s = "hello world"
```
