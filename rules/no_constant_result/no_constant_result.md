---
title: noConstantResult
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noConstantResult`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noConstantResult:
    check-exported: false # Whether to check exported (public) functions. Default: false.
```

## Details

Reports named result parameters (return values) that always return the same value across all execution paths. When a function's named return value is invariably set to the same constant (such as always returning `nil` for an error), the result parameter adds unnecessary complexity to the function signature.

This check uses whole-program SSA analysis to trace all return paths and determine if a named result parameter is always assigned the same value. If so, the parameter can typically be removed from the return signature, and callers can be updated to no longer expect that return value.

A common occurrence of this pattern is an error return value that is always `nil`, indicating the function never actually fails. This misleads callers into adding error-handling code for a condition that can never occur, adding unnecessary complexity. Removing such a return value simplifies both the function and its call sites.

By default, this rule only checks unexported (private) functions. Set `check-exported: true` to also analyze exported functions.

The typical diagnostic message is: `result resultName is always nil` (or another constant value).

Source: https://github.com/mvdan/unparam

## Examples

### Invalid

```golang
// The error return value is always nil.
func parseConfig(path string) (Config, error) {
    data := defaultConfig()
    return data, nil
}

func main() {
    cfg, err := parseConfig("/etc/app.conf")
    if err != nil {
        // This branch can never be reached.
        log.Fatal(err)
    }
    use(cfg)
}
```

```golang
// The boolean return value is always true.
func validate(name string) (string, bool) {
    cleaned := strings.TrimSpace(name)
    return cleaned, true
}
```

```golang
// The second return value "count" is always 0.
func process(items []string) (result []string, count int) {
    for _, item := range items {
        result = append(result, strings.ToUpper(item))
    }
    return result, 0
}
```

### Valid

```golang
// The error return value varies depending on the execution path.
func parseConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("reading config: %w", err)
    }
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return Config{}, fmt.Errorf("parsing config: %w", err)
    }
    return cfg, nil
}
```

```golang
// The boolean return value varies depending on the input.
func validate(name string) (string, bool) {
    cleaned := strings.TrimSpace(name)
    if cleaned == "" {
        return "", false
    }
    return cleaned, true
}
```

```golang
// The count return value reflects actual processing.
func process(items []string) (result []string, count int) {
    for _, item := range items {
        if item != "" {
            result = append(result, strings.ToUpper(item))
            count++
        }
    }
    return result, count
}
```
