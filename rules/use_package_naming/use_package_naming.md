---
title: usePackageNaming
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/usePackageNaming`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/usePackageNaming:
    arguments:
      - skipConventionNameCheck: false
        conventionNameCheckRegex: ""
        skipTopLevelCheck: false
        skipDefaultBadNameCheck: false
        checkExtraBadName: false
        userDefinedBadNames: []
        skipCollisionWithCommonStd: false
        checkCollisionWithAllStd: false
```

## Details

This rule checks that package names follow [Go conventions](https://go.dev/blog/package-names) and best practices.
It helps prevent using bad package names and enforces consistent naming patterns.
This rule arose from package naming checks in `var-naming`.

By default, it checks for:

- Package name conventions (no underscores except for test packages, no MixedCaps).
- Bad package names from the official Go blog (e.g., `common`, `util`, `utils`, `misc`, `interfaces`, `types`).
- Package names that conflict with common Go standard library packages (e.g., `http`, `json`, `fmt`).

**Configuration options** (optional single map argument):

- `skipConventionNameCheck` (bool): If `true`, skip checks for package name conventions (underscores, MixedCaps, etc.). Default: `false`. Mutually exclusive with `conventionNameCheckRegex`.
- `conventionNameCheckRegex` (string): Custom regex pattern to validate package names. If set, package names must match this pattern. Mutually exclusive with `skipConventionNameCheck`.
- `skipTopLevelCheck` (bool): If `true`, skip checks for top-level package names (e.g., `pkg`). Default: `false`.
- `skipDefaultBadNameCheck` (bool): If `true`, skip checks for default bad package names (e.g., `common`, `utils`). Default: `false`.
- `checkExtraBadName` (bool): If `true`, enable checks for extra bad package names (e.g., `helpers`, `models`, `shared`, `utilities`). Default: `false`.
- `userDefinedBadNames` ([]string): List of user-defined bad package names to check for.
- `skipCollisionWithCommonStd` (bool): If `true`, skip checks for collisions with the most common Go standard library packages. Default: `false`.
- `checkCollisionWithAllStd` (bool): If `true`, enable checks for collisions with all packages from Go standard library. Default: `false`. Mutually exclusive with `skipCollisionWithCommonStd`.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
package PkgName // don't use package name that contains MixedCaps
```

```golang
package pkg_name // don't use package name that contains an underscore
```

```golang
package util // don't use "util" because it is a bad package name
```

```golang
package http // don't use "http" because it conflicts with common Go standard library package "net/http"
```

### Valid

```golang
package mypackage
```

```golang
package server
```
