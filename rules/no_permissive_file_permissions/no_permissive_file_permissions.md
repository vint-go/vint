---
title: noPermissiveFilePermissions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPermissiveFilePermissions`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPermissiveFilePermissions:
    maxPermission: "0600"
```

## Details

Detects poor file permissions used when creating a file or using `chmod`.

The `maxPermission` option allows customizing the maximum allowed file permission mode. The value can be specified as an octal string (e.g., `"0600"`) or decimal integer. Default is `0600`.

This rule monitors calls to `os.OpenFile` and `os.Chmod` and flags cases where the permission mode is more permissive than `0600` (default threshold). Overly permissive file permissions can allow other users on the system to read or modify files that should be private.

The default maximum permission is `0600`, which grants only the file owner read and write access. Files with group or world permissions (e.g., `0644`, `0666`, `0777`) may expose sensitive data in multi-user environments.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func createFile() (*os.File, error) {
    // World-readable file permissions
    return os.OpenFile("config.yml", os.O_CREATE|os.O_WRONLY, 0644)
}
```

```golang
import "os"

func setPermissions() error {
    // Overly permissive chmod
    return os.Chmod("secret.key", 0666)
}
```

```golang
import "os"

func createExecutable() (*os.File, error) {
    // World-executable file
    return os.OpenFile("script.sh", os.O_CREATE|os.O_WRONLY, 0755)
}
```

### Valid

```golang
import "os"

func createFile() (*os.File, error) {
    // Owner-only read/write permissions
    return os.OpenFile("config.yml", os.O_CREATE|os.O_WRONLY, 0600)
}
```

```golang
import "os"

func setPermissions() error {
    // Owner-only permissions
    return os.Chmod("secret.key", 0600)
}
```
