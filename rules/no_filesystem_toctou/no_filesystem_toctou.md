---
title: noFilesystemToctou
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noFilesystemToctou`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noFilesystemToctou:
    # rule options here
```

## Details

Detects filesystem TOCTOU (Time-of-Check-Time-of-Use) race conditions in `filepath.Walk` and `filepath.WalkDir` callbacks.

This rule uses SSA analysis to identify patterns where a file's properties are checked during the walk callback but the file is then accessed in a way that is vulnerable to race conditions. Between the time the callback receives file information and the time the file is actually accessed, an attacker could replace the file with a symlink or different file, leading to unauthorized file access, path traversal, or privilege escalation.

TOCTOU vulnerabilities in filesystem operations are particularly dangerous in setuid programs or services that process files in shared or user-controlled directories.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "os"
    "path/filepath"
)

func processFiles(root string) error {
    return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            // TOCTOU: file could change between Walk's stat and this Open
            data, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            process(data)
        }
        return nil
    })
}
```

### Valid

```golang
import (
    "os"
    "path/filepath"
)

func processFiles(root string) error {
    return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            // Open the file and use Fstat on the file descriptor
            f, err := os.Open(path)
            if err != nil {
                return err
            }
            defer f.Close()

            fi, err := f.Stat()
            if err != nil {
                return err
            }
            // Verify the file is still a regular file
            if !fi.Mode().IsRegular() {
                return nil
            }
            data, err := io.ReadAll(f)
            if err != nil {
                return err
            }
            process(data)
        }
        return nil
    })
}
```
