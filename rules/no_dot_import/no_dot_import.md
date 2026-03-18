---
title: noDotImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noDotImport`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noDotImport:
    # rule options here
```

## Details

Dot imports are discouraged.

Dot imports (`import . "pkg"`) make code harder to read because it becomes unclear which package a name belongs to. Use regular imports instead.

Source: https://staticcheck.dev/docs/checks/#ST1001

## Examples

### Invalid

```golang
package main

import . "fmt"

func main() {
    // Unclear that Println comes from fmt
    Println("hello")
}
```

### Valid

```golang
package main

import "fmt"

func main() {
    fmt.Println("hello")
}
```
