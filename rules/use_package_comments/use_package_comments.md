---
title: usePackageComments
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/usePackageComments`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/usePackageComments:
    # no configuration options
```

## Details

Packages should have comments. This rule warns on undocumented packages and when package comments are detached from the `package` keyword.

It checks that:
- The package has a doc comment on at least one file.
- The package comment is not detached (no blank lines between the comment and the `package` statement).
- The package comment follows the standard format `Package <name> ...` for non-main packages.

More information: https://go.dev/wiki/CodeReviewComments#package-comments

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// This package does stuff
package mypackage
```

```golang
// Package mypackage provides utilities.

package mypackage // detached comment
```

### Valid

```golang
// Package mypackage provides utilities for processing data.
package mypackage
```
