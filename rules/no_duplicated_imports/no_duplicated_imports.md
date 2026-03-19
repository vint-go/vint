---
title: noDuplicatedImports
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noDuplicatedImports`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noDuplicatedImports:
    # no additional options
```

## Details

It is possible to unintentionally import the same package twice. This rule looks for packages that are imported two or more times.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
import (
    "crypto/md5"
    _ "crypto/md5"
)
```

```golang
import (
    "strings"
    str "strings"
)
```

### Valid

```golang
import (
    "crypto/md5"
    "strings"
)
```
