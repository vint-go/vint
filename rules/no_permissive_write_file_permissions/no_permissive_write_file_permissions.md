---
title: noPermissiveWriteFilePermissions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPermissiveWriteFilePermissions`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPermissiveWriteFilePermissions:
    # rule options here
```

## Details

Detects poor file permissions used when writing to a file with `os.WriteFile` or `ioutil.WriteFile`.

This rule flags calls to `os.WriteFile` (and the deprecated `ioutil.WriteFile`) where the permission mode is more permissive than `0600` (default threshold). When writing sensitive data to files, overly permissive permissions can allow other users on the system to read the file contents.

The default maximum permission is `0600`, which grants only the file owner read and write access. Files containing configuration data, credentials, tokens, or other sensitive information should use restrictive permissions.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func saveConfig(data []byte) error {
    // World-readable file permissions
    return os.WriteFile("config.yml", data, 0644)
}
```

```golang
import "os"

func saveKey(key []byte) error {
    // Overly permissive for sensitive data
    return os.WriteFile("private.key", key, 0666)
}
```

### Valid

```golang
import "os"

func saveConfig(data []byte) error {
    // Owner-only read/write permissions
    return os.WriteFile("config.yml", data, 0600)
}
```

```golang
import "os"

func saveKey(key []byte) error {
    // Owner-only read permissions for key file
    return os.WriteFile("private.key", key, 0400)
}
```
