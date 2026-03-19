---
title: noBlankImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noBlankImport`
- This rule is recommended, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noBlankImport:
    # no configuration options
```

## Details

Blank import should be only in a main or test package, or have a comment justifying it.

A blank import (e.g. `_ "some/package"`) causes the imported package's `init()` function to run for its side effects. In library packages, blank imports without an explanatory comment can be confusing and make the code harder to maintain. This rule requires that the first blank import in each contiguous group either appears in a `main` or `_test` package, or has a comment (doc or inline) explaining why the side-effect import is needed.

The rule also recognizes the `embed` package: a blank import of `"embed"` is allowed without a comment when the file contains a valid `//go:embed` directive.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// In a library package (not main, not test):
package mylib

import (
    _ "image/png"  // no justifying comment - flagged
)
```

```golang
// In a library package with multiple blank imports:
package mylib

import (
    _ "image/gif"  // no justifying comment on first blank import - flagged
    _ "image/jpeg" // subsequent in same group - not flagged
)
```

### Valid

```golang
// In a main package - blank imports are always allowed:
package main

import _ "image/png"
```

```golang
// In a library package with justifying comments:
package mylib

import (
    // PNG decoding support
    _ "image/png"

    _ "image/jpeg" // JPEG decoding support
)
```

```golang
// Blank import of embed with a valid go:embed directive:
package mylib

import _ "embed"

//go:embed template.html
var tmpl string
```
