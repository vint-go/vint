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
    arguments:
      - allowedPackages:
          - "github.com/onsi/ginkgo/v2"
          - "github.com/onsi/gomega"
```

## Details

Importing with `.` makes the programs much harder to understand because it is unclear whether names belong to the current package or to an imported package.

More information: https://go.dev/wiki/CodeReviewComments#import-dot

The `allowedPackages` option lets you whitelist specific packages that are permitted to use dot imports (e.g. testing DSL packages like Ginkgo/Gomega).

Source: https://github.com/mgechev/revive

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
