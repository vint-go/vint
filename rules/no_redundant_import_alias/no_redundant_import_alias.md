---
title: noRedundantImportAlias
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRedundantImportAlias`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRedundantImportAlias:
    # no options
```

## Details

This rule warns on redundant import aliases. This happens when the alias used on the import statement matches the imported package name. Such aliases add visual noise without providing any benefit.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
import (
	md5 "crypto/md5"       // alias "md5" matches package name
	strings "strings"      // alias "strings" matches package name
)
```

### Valid

```golang
import (
	"crypto/md5"           // no alias
	_ "crypto/md5"         // blank import
	crypto "crypto/md5"    // alias differs from package name
	str "strings"          // alias differs from package name
)
```
