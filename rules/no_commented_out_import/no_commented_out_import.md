---
title: noCommentedOutImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noCommentedOutImport`
- This rule is **not recommended**, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noCommentedOutImport:
    # no additional options
```

## Details

Detects commented-out imports. This checker scans comments inside import declaration blocks to find commented-out import paths and suggests removing them. Commented-out imports are dead code and should be removed or managed via version control.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
import (
    "fmt"
    // "os"
    "strings"
)
```

### Valid

```golang
import (
    "fmt"
    "strings"
)
```
