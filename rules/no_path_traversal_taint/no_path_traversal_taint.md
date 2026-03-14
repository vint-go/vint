---
title: noPathTraversalTaint
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/security/noPathTraversalTaint`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/security/noPathTraversalTaint:
    # rule options here
```

## Details

Detects path traversal vulnerabilities via taint analysis.

This rule uses taint analysis to trace the flow of untrusted data from sources to filesystem access sinks. It performs interprocedural data flow analysis to detect cases where user-controlled file paths flow through multiple functions before reaching file operations such as `os.Open`, `os.ReadFile`, `os.Create`, or `os.WriteFile`.

Path traversal attacks (also known as directory traversal) allow attackers to access files outside the intended directory by using sequences like `../` or absolute paths. This can lead to reading sensitive configuration files, overwriting critical system files, or accessing other users' data.

File paths derived from user input should be validated by resolving the absolute path and confirming it stays within the intended directory. Using `filepath.Clean` alone is insufficient; the resolved path must be checked against the allowed base directory.

Source: https://github.com/securego/gosec

## Examples

### Invalid

```golang
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    // Taint flows from request to file read
    data := readUserFile(filename)
    w.Write(data)
}

func readUserFile(name string) []byte {
    path := filepath.Join("/uploads", name)
    data, _ := os.ReadFile(path) // Path traversal possible
    return data
}
```

### Valid

```golang
func downloadHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    data, err := readUserFileSafe(filename)
    if err != nil {
        http.Error(w, "not found", http.StatusNotFound)
        return
    }
    w.Write(data)
}

func readUserFileSafe(name string) ([]byte, error) {
    baseDir := "/uploads"
    // Resolve the absolute path
    absPath, err := filepath.Abs(filepath.Join(baseDir, name))
    if err != nil {
        return nil, err
    }
    // Verify the path stays within the base directory
    if !strings.HasPrefix(absPath, baseDir+string(os.PathSeparator)) {
        return nil, fmt.Errorf("path traversal detected")
    }
    return os.ReadFile(absPath)
}
```

```golang
// Go 1.24+: Using os.Root for safe scoped access
func readUserFileSafe(name string) ([]byte, error) {
    root, err := os.OpenRoot("/uploads")
    if err != nil {
        return nil, err
    }
    defer root.Close()
    f, err := root.Open(name)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return io.ReadAll(f)
}
```
