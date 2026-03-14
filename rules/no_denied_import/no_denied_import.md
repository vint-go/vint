---
title: noDeniedImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noDeniedImport`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noDeniedImport:
    rules:
      main:
        list-mode: lax
        deny:
          - pkg: "github.com/sirupsen/logrus"
            desc: "Use log/slog instead of logrus"
          - pkg: "github.com/pkg/errors"
            desc: "Use fmt.Errorf with %%w verb for error wrapping"
```

## Details

Reports when a Go source file imports a package that appears on the deny list configured for the matching depguard rule group.

Depguard is a Go linter that controls which packages are allowed to be imported. The `noDeniedImport` check triggers when a package matches an entry in the `deny` list of any applicable rule group. Each deny entry can include a descriptive suggestion message that will be shown alongside the diagnostic, guiding the developer toward an acceptable alternative.

Package matching uses prefix matching by default. For example, denying `github.com/foo/bar` will also deny `github.com/foo/bar/baz`. To require an exact match (no prefix matching), append `$` to the package name (e.g., `github.com/foo/bar$`).

Rules are organized into named groups, and each group can specify which files it applies to using glob patterns in the `files` field. The special variable `$all` matches all Go files, and `$test` matches only test files. File patterns can be negated with a `!` prefix.

When `list-mode` is set to `lax` (the default-allow model), all packages are allowed unless explicitly listed in the `deny` section. When set to `strict`, both `allow` and `deny` interact: a denied package is blocked unless a more specific allow entry overrides it.

Source: https://github.com/OpenPeeDeeP/depguard

## Examples

### Invalid

```golang
// .depguard.yaml denies "github.com/sirupsen/logrus"
package main

import (
    "github.com/sirupsen/logrus" // error: import 'github.com/sirupsen/logrus' is not allowed from list 'main': Use log/slog instead of logrus
)

func main() {
    logrus.Info("hello")
}
```

```golang
// .depguard.yaml denies "github.com/pkg/errors"
package main

import (
    "github.com/pkg/errors" // error: import 'github.com/pkg/errors' is not allowed from list 'main': Use fmt.Errorf with %w verb for error wrapping
)

func main() {
    err := errors.New("something failed")
    _ = err
}
```

```golang
// .depguard.yaml denies "io/ioutil" (deprecated package)
package main

import (
    "io/ioutil" // error: import 'io/ioutil' is not allowed from list 'main': Use os and io packages directly instead
)

func main() {
    data, _ := ioutil.ReadFile("test.txt")
    _ = data
}
```

### Valid

```golang
// Using an allowed alternative instead of a denied package
package main

import (
    "log/slog"
)

func main() {
    slog.Info("hello")
}
```

```golang
// Using fmt.Errorf instead of denied github.com/pkg/errors
package main

import (
    "fmt"
)

func main() {
    err := fmt.Errorf("something failed: %w", fmt.Errorf("cause"))
    _ = err
}
```

```golang
// Using os.ReadFile instead of denied io/ioutil
package main

import (
    "os"
)

func main() {
    data, _ := os.ReadFile("test.txt")
    _ = data
}
```
