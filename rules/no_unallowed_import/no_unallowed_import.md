---
title: noUnallowedImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUnallowedImport`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUnallowedImport:
    rules:
      main:
        list-mode: strict
        files:
          - $all
        allow:
          - $gostd
          - github.com/myorg
```

## Details

Reports when a Go source file imports a package that is not present in the allow list configured for the matching depguard rule group.

This check is particularly relevant when using the `strict` list mode. In strict mode, every import must be explicitly allowed -- any package that does not match an entry in the `allow` list is denied by default. This is useful for teams that want tight control over their dependency graph and want to ensure only vetted packages are used.

The `allow` list uses prefix matching by default. For example, allowing `github.com/myorg` will also allow `github.com/myorg/pkg1`, `github.com/myorg/pkg2`, and any other sub-package. To require exact matching, append `$` to the package name (e.g., `github.com/myorg/pkg1$` allows only that exact package, not its sub-packages).

The special variable `$gostd` can be used in the allow list to permit all Go standard library packages.

In `original` list mode (the legacy default), if an allow list is defined, packages not matching it are also denied. In `lax` mode, packages are allowed by default and only the deny list is enforced, making this check less relevant.

Rules are organized into named groups, each targeting specific files via glob patterns. The `$all` variable matches all Go files, while `$test` matches only test files. File patterns can be negated with a `!` prefix.

Source: https://github.com/OpenPeeDeeP/depguard

## Examples

### Invalid

```golang
// strict mode: only $gostd and github.com/myorg are allowed
package main

import (
    "fmt"
    "github.com/some/other-lib" // error: import 'github.com/some/other-lib' is not allowed from list 'main'
)

func main() {
    fmt.Println("hello")
}
```

```golang
// strict mode: only specific packages are allowed
package main

import (
    "github.com/myorg/approved-pkg"
    "github.com/thirdparty/unapproved" // error: import 'github.com/thirdparty/unapproved' is not allowed from list 'main'
)

func main() {
    _ = approved-pkg.New()
}
```

```golang
// strict mode with exact matching: github.com/myorg/utils$ is allowed but not sub-packages
package main

import (
    "github.com/myorg/utils/internal" // error: import 'github.com/myorg/utils/internal' is not allowed from list 'main'
)

func main() {
}
```

### Valid

```golang
// All imports are on the allow list
package main

import (
    "fmt"
    "net/http"
    "github.com/myorg/service"
)

func main() {
    fmt.Println("hello")
    _ = http.StatusOK
    _ = service.New()
}
```

```golang
// Standard library packages are allowed via $gostd
package main

import (
    "context"
    "encoding/json"
    "os"
)

func main() {
    ctx := context.Background()
    _ = ctx
    data, _ := os.ReadFile("config.json")
    _ = json.Unmarshal(data, nil)
}
```

```golang
// Organizational packages allowed via prefix matching
package main

import (
    "github.com/myorg/auth"
    "github.com/myorg/database"
    "github.com/myorg/logging"
)

func main() {
    _ = auth.NewClient()
    _ = database.Connect()
    logging.Info("started")
}
```
