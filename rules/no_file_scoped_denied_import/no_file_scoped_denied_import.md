---
title: noFileScopedDeniedImport
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noFileScopedDeniedImport`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noFileScopedDeniedImport:
    rules:
      main:
        files:
          - "!$test"
        deny:
          - pkg: "github.com/stretchr/testify"
            desc: "testify should only be used in test files"
      tests:
        files:
          - $test
        allow:
          - $gostd
          - github.com/stretchr/testify
```

## Details

Reports when a package import violates a file-scoped depguard rule, meaning the import is prohibited specifically in certain types of files based on glob patterns.

Depguard allows defining multiple named rule groups, each targeting a specific set of files via the `files` field. This enables fine-grained control over which packages are permitted in which parts of the codebase. For example, you can deny test-only packages like `github.com/stretchr/testify` in production code while allowing them in test files.

The `files` field accepts glob patterns and the following built-in variables:

- `$all` - Matches all Go files.
- `$test` - Matches only test files (files ending with `_test.go`).
- `!` prefix - Negates a pattern. For example, `!$test` matches all non-test files.

When multiple rule groups are defined, a file may be subject to more than one group if its path matches the file patterns of multiple groups. Each matching group's allow and deny lists are evaluated independently.

This mechanism is commonly used to enforce architectural boundaries, such as:
- Preventing test dependencies from leaking into production code.
- Restricting internal packages from being imported in public-facing code.
- Applying different dependency policies to different modules or directories.

Source: https://github.com/OpenPeeDeeP/depguard

## Examples

### Invalid

```golang
// Rule: testify is denied in non-test files (files: ["!$test"])
// File: handler.go
package handler

import (
    "net/http"
    "github.com/stretchr/testify/assert" // error: import 'github.com/stretchr/testify/assert' is not allowed from list 'main': testify should only be used in test files
)

func Handler(w http.ResponseWriter, r *http.Request) {
    // Production code should not use testify
}
```

```golang
// Rule: internal packages denied outside their module (files: ["!**/internal/**"])
// File: cmd/server/main.go
package main

import (
    "github.com/myorg/project/internal/secret" // error: import 'github.com/myorg/project/internal/secret' is not allowed from list 'restricted'
)

func main() {
    _ = secret.GetKey()
}
```

```golang
// Rule: database packages denied in handler files (files: ["**/handler/**"])
// File: handler/user.go
package handler

import (
    "database/sql" // error: import 'database/sql' is not allowed from list 'handlers': Use the repository layer for database access
)

func GetUser(id int) {
    // Handlers should not directly access the database
}
```

### Valid

```golang
// testify is allowed in test files
// File: handler_test.go
package handler_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    assert.NotNil(t, req)
}
```

```golang
// Production code uses only allowed packages
// File: handler.go
package handler

import (
    "encoding/json"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

```golang
// Handler uses repository layer instead of direct database access
// File: handler/user.go
package handler

import (
    "net/http"
    "github.com/myorg/project/repository"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
    user, _ := repository.FindUser(1)
    _ = user
}
```
