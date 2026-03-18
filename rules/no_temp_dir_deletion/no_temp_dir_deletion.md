---
title: noTempDirDeletion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noTempDirDeletion`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noTempDirDeletion:
    # rule options here
```

## Details

Deleting a directory that is prepopulated with files from `os.TempDir`.

Calling `os.RemoveAll(os.TempDir())` or similar operations on the system temp directory can delete files belonging to other processes. Create a subdirectory within the temp directory instead.

Source: https://staticcheck.dev/docs/checks/#SA9007

## Examples

### Invalid

```golang
package main

import "os"

func cleanup() {
    // Dangerous: deletes the entire temp directory
    os.RemoveAll(os.TempDir())
}
```

### Valid

```golang
package main

import "os"

func cleanup() {
    dir, _ := os.MkdirTemp("", "myapp-")
    defer os.RemoveAll(dir)
    // use dir...
}
```
