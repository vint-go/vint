---
title: noPredictableTempFile
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPredictableTempFile`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPredictableTempFile:
    # rule options here
```

## Details

Detects the creation of temporary files using a predictable path.

This rule flags calls that create files in temporary directories using hardcoded, predictable file names instead of using functions that generate unique file names (like `os.CreateTemp` or `ioutil.TempFile`). Predictable temporary file paths are vulnerable to symlink attacks where an attacker creates a symbolic link at the expected path before the application does, potentially causing the application to write to or read from an attacker-controlled location.

Applications should use `os.CreateTemp` (Go 1.16+) or `ioutil.TempFile` (deprecated) to create temporary files with unique, unpredictable names. These functions use atomic operations to prevent race conditions.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func writeTemp() error {
    // Predictable temporary file path
    f, err := os.Create("/tmp/myapp.log")
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.WriteString("log data")
    return err
}
```

```golang
import "os"

func writeTempData(data []byte) error {
    // Predictable path in temp directory
    return os.WriteFile("/tmp/myapp-cache.dat", data, 0600)
}
```

### Valid

```golang
import "os"

func writeTemp() error {
    // Using os.CreateTemp for unique file name
    f, err := os.CreateTemp("", "myapp-*.log")
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.WriteString("log data")
    return err
}
```

```golang
import "os"

func writeTempData(data []byte) error {
    // Using os.CreateTemp with a specific directory
    f, err := os.CreateTemp(os.TempDir(), "myapp-cache-*.dat")
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = f.Write(data)
    return err
}
```
