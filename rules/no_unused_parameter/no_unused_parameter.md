---
title: noUnusedParameter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noUnusedParameter`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noUnusedParameter:
    check-exported: false # Whether to check exported (public) functions. Default: false.
```

## Details

Reports function parameters that are completely unused within the function body. Unlike the compiler's unused variable check, `unparam` uses whole-program analysis via the SSA (Static Single Assignment) intermediate representation to detect parameters that serve no purpose in the function's logic.

Unused parameters increase code complexity, mislead readers about the function's dependencies, and may indicate incomplete refactoring. They add cognitive overhead for anyone reading or maintaining the code, since the reader must determine why the parameter exists and whether it should be used.

If a parameter is genuinely needed to satisfy an interface or a function signature requirement (e.g., `http.HandlerFunc`), consider using the blank identifier `_` to explicitly mark it as intentionally unused.

By default, this rule only checks unexported (private) functions. Exported functions are excluded because they may need to satisfy interface contracts or maintain API compatibility. Set `check-exported: true` to also analyze exported functions.

Source: https://github.com/mvdan/unparam

## Examples

### Invalid

```golang
// The parameter "name" is never used in the function body.
func greet(name string) string {
    return "hello"
}
```

```golang
// The parameter "count" is never referenced.
func processItems(items []string, count int) error {
    for _, item := range items {
        fmt.Println(item)
    }
    return nil
}
```

```golang
// The parameter "logger" is unused; the function uses fmt directly.
func fetchData(ctx context.Context, url string, logger *log.Logger) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

### Valid

```golang
// All parameters are used in the function body.
func greet(name string) string {
    return "hello, " + name
}
```

```golang
// All parameters are referenced.
func processItems(items []string, count int) error {
    if len(items) != count {
        return fmt.Errorf("expected %d items, got %d", count, len(items))
    }
    for _, item := range items {
        fmt.Println(item)
    }
    return nil
}
```

```golang
// Using the blank identifier to explicitly mark an unused parameter
// that is required by an interface or callback signature.
func handler(_ http.ResponseWriter, r *http.Request) {
    log.Println("received request:", r.URL.Path)
}
```

```golang
// All parameters are used, including the logger.
func fetchData(ctx context.Context, url string, logger *log.Logger) ([]byte, error) {
    logger.Printf("fetching %s", url)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```
