---
title: noFileInclusionViaVariable
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noFileInclusionViaVariable`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noFileInclusionViaVariable:
    # rule options here
```

## Details

Detects potential file inclusion via variable, where file paths provided as taint input are used to open or read files.

This rule monitors calls to `os.Open`, `os.OpenFile`, `os.ReadFile`, `os.Create`, and `ioutil.ReadFile`, and flags cases where the file path argument is derived from a variable rather than a constant. Variable file paths can be manipulated by attackers to perform directory traversal attacks, accessing files outside the intended directory.

The rule recognizes path sanitization functions such as `filepath.Clean`, `filepath.Rel`, and `filepath.EvalSymlinks`, and considers paths that pass through these functions as safer. For Go 1.24+, the rule suggests using `os.Root` to scope file access under a fixed root directory, preventing directory traversal.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
import "os"

func readFile(path string) ([]byte, error) {
    // File path from variable - potential directory traversal
    return os.ReadFile(path)
}
```

```golang
import "os"

func openFile(userPath string) (*os.File, error) {
    // User-controlled path
    return os.Open(userPath)
}
```

```golang
import "os"

func readUserFile(filename string) ([]byte, error) {
    // Concatenation with user input
    path := "/data/uploads/" + filename
    return os.ReadFile(path)
}
```

### Valid

```golang
import (
    "os"
    "path/filepath"
)

func readFile(basePath, filename string) ([]byte, error) {
    // Clean and validate the path
    cleanPath := filepath.Clean(filepath.Join(basePath, filename))
    if !strings.HasPrefix(cleanPath, basePath) {
        return nil, fmt.Errorf("path traversal detected")
    }
    return os.ReadFile(cleanPath)
}
```

```golang
import "os"

func readConfig() ([]byte, error) {
    // Constant file path
    return os.ReadFile("/etc/myapp/config.yaml")
}
```

```golang
// Go 1.24+: Using os.Root for scoped file access
import "os"

func readScopedFile(filename string) ([]byte, error) {
    root, err := os.OpenRoot("/data/uploads")
    if err != nil {
        return nil, err
    }
    defer root.Close()
    f, err := root.Open(filename)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return io.ReadAll(f)
}
```
