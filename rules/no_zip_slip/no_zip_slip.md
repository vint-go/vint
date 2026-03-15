---
title: noZipSlip
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noZipSlip`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noZipSlip:
    # rule options here
```

## Details

Detects file path traversal when extracting zip or tar archives (Zip Slip vulnerability).

This rule identifies code that extracts archive entries without validating that the resulting file path stays within the intended destination directory. Maliciously crafted archive files can contain entries with paths like `../../etc/cron.d/malicious` that, when extracted, write files outside the intended extraction directory. This is known as the "Zip Slip" vulnerability.

When extracting archives, the code should validate that each entry's resolved path starts with the intended destination directory. Using `filepath.Join` alone is not sufficient; the resulting path must be checked to ensure it does not escape the target directory.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import (
    "archive/zip"
    "io"
    "os"
    "path/filepath"
)

func extractZip(zipPath, dest string) error {
    r, _ := zip.OpenReader(zipPath)
    defer r.Close()

    for _, f := range r.File {
        // Path traversal: entry name could be "../../etc/passwd"
        path := filepath.Join(dest, f.Name)
        outFile, _ := os.Create(path)
        rc, _ := f.Open()
        io.Copy(outFile, rc)
        outFile.Close()
        rc.Close()
    }
    return nil
}
```

### Valid

```golang
import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
)

func extractZip(zipPath, dest string) error {
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return err
    }
    defer r.Close()

    dest = filepath.Clean(dest) + string(os.PathSeparator)

    for _, f := range r.File {
        path := filepath.Join(dest, f.Name)
        // Validate that the path stays within the destination
        if !strings.HasPrefix(path, dest) {
            return fmt.Errorf("illegal file path: %s", f.Name)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, 0750)
            continue
        }

        os.MkdirAll(filepath.Dir(path), 0750)
        outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
            return err
        }

        rc, err := f.Open()
        if err != nil {
            outFile.Close()
            return err
        }

        _, err = io.Copy(outFile, rc)
        outFile.Close()
        rc.Close()
        if err != nil {
            return err
        }
    }
    return nil
}
```
