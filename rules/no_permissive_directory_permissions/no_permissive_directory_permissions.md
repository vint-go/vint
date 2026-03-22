---
title: noPermissiveDirectoryPermissions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPermissiveDirectoryPermissions`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPermissiveDirectoryPermissions:
    maxPermission: "0750"
```

## Details

Detects poor file permissions used when creating a directory.

The `maxPermission` option allows customizing the maximum allowed directory permission mode. The value can be specified as an octal string (e.g., `"0750"`) or decimal integer. Default is `0750`.

This rule monitors calls to `os.Mkdir` and `os.MkdirAll` and flags cases where the permission mode is more permissive than `0750` (default threshold). Overly permissive directory permissions can allow other users on the system to read, write, or traverse directories that should be restricted.

The default maximum permission is `0750`, which grants the owner full access and the group read and execute access. Directories with world-readable or world-writable permissions (e.g., `0777`) are a security risk in multi-user environments.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func createDir() error {
    // Overly permissive directory permissions
    return os.Mkdir("/tmp/myapp", 0777)
}
```

```golang
import "os"

func createDirs() error {
    // World-readable directory tree
    return os.MkdirAll("/var/myapp/data", 0755)
}
```

### Valid

```golang
import "os"

func createDir() error {
    // Restrictive directory permissions
    return os.Mkdir("/tmp/myapp", 0750)
}
```

```golang
import "os"

func createDirs() error {
    // Owner-only directory permissions
    return os.MkdirAll("/var/myapp/data", 0700)
}
```
