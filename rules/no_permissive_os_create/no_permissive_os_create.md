---
title: noPermissiveOsCreate
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPermissiveOsCreate`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPermissiveOsCreate:
    # rule options here
```

## Details

Detects poor file permissions used when creating a file with `os.Create`.

The `os.Create` function creates a file with default permissions of `0666` (before umask), which can result in files that are readable and writable by all users on the system. This is particularly concerning when the file will contain sensitive data such as configuration, credentials, or private keys.

Instead of using `os.Create`, consider using `os.OpenFile` with explicit, restrictive permissions (e.g., `0600`) to ensure the file is only accessible by the owner.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func createFile() (*os.File, error) {
    // os.Create uses default 0666 permissions
    return os.Create("sensitive-data.txt")
}
```

```golang
import "os"

func writeSecret(secret string) error {
    f, err := os.Create("secret.key") // Permissive default permissions
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.WriteString(secret)
    return err
}
```

### Valid

```golang
import "os"

func createFile() (*os.File, error) {
    // Explicit restrictive permissions
    return os.OpenFile("sensitive-data.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}
```

```golang
import "os"

func writeSecret(secret string) error {
    return os.WriteFile("secret.key", []byte(secret), 0600)
}
```
