---
title: useFileHeader
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useFileHeader`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useFileHeader:
    arguments:
      - "This is the text that must appear at the top of source files."
```

## Details

This rule helps to enforce a common header for all source files in a project by spotting those files that do not have the specified header. The argument is a regular expression that the file header comment must match.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// This file has no matching header
package main

func main() {}
```

### Valid

```golang
// This is the text that must appear at the top of source files.

package main

func main() {}
```
