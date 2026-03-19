---
title: noPackageDirectoryMismatch
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noPackageDirectoryMismatch`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noPackageDirectoryMismatch:
    arguments:
      - ignoreDirectories: ["testcases", "testinfo"]
```

## Details

It is considered a good practice to name a package after the directory containing it.
This rule warns when the package name declared in the file does not match the name of the directory containing the file.

The following cases are excluded from this check:

- Package `main` (executable packages)
- Files in `testdata` directories (at any level) - by default
- Files directly in `internal` directories (but files in subdirectories of `internal` are checked)
- Root directories (directories containing `go.mod` or `.git`)

For test files (files with `_test` suffix), the package name is additionally checked to see if it matches the directory name with `_test` suffix appended. The special package name `main_test` is always allowed in test files.

The rule normalizes both directory and package names before comparison by removing hyphens (`-`),
underscores (`_`), and dots (`.`). This allows package `foo_barbuz` to be considered equal to directory `foo-bar.buz`.

For files in version directories (`v1`, `v2`, etc.), the package name is checked against both the version directory and its parent directory.

The `ignoreDirectories` option allows specifying a list of directory name patterns to exclude from the check.
By default, directories matching `testdata` are ignored. Passing an empty list enables checking all directories.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// File: mylib/foo.go
package bar // package name "bar" does not match directory name "mylib"
```

```golang
// File: api/v1v/handler.go
package api // package name "api" does not match directory name "v1v"
```

### Valid

```golang
// File: mylib/foo.go
package mylib // package name matches directory name
```

```golang
// File: go-mylib/foo.go
package mylib // "go-" prefix is stripped during normalization
```

```golang
// File: api/v1/handler.go
package api // version directory v1, parent directory "api" matches
```

```golang
// File: cmd/server/main.go
package main // main packages are always allowed
```

```golang
// File: mylib/foo_test.go
package mylib_test // external test package matches directory + "_test"
```
