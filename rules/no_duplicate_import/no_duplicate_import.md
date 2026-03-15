---
title: noDuplicateImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noDuplicateImport`
- This rule is not recommended (experimental).
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noDuplicateImport:
    # no additional options
```

## Details

Detects multiple imports of the same package under different aliases. Importing the same package multiple times with different aliases creates confusion and should be consolidated into a single import.

Source: https://github.com/go-critic/go-critic

## Examples

### Invalid

```golang
import (
    "fmt"
    printing "fmt"
)
```

### Valid

```golang
import (
    "fmt"
)
```
