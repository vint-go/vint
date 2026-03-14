---
title: noLongFunctions
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noLongFunctions`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noLongFunctions:
    lines: 60
    ignoreComments: false
```

## Details

Checks that functions do not exceed a maximum number of lines. Long functions are harder to read, understand, and maintain. They often indicate that a function is doing too much and should be broken down into smaller, more focused functions.

The default maximum is 60 lines per function. This limit is intended to ensure that a function fits within a single screen, making it easier to trace variable definitions and find matching brackets without scrolling.

Lines are counted from the opening brace to the closing brace of the function body, excluding the function signature itself. When `ignoreComments` is set to `true`, comment lines inside the function body are excluded from the line count.

The line limit can be set to `-1` to disable this check entirely.

Source: https://github.com/ultraware/funlen

## Examples

### Invalid

```golang
// This function exceeds the default 60-line limit.
func processData(input []string) ([]string, error) {
    result := make([]string, 0, len(input))
    errors := make([]error, 0)

    for _, item := range input {
        if item == "" {
            continue
        }

        trimmed := strings.TrimSpace(item)
        if len(trimmed) == 0 {
            continue
        }

        if strings.HasPrefix(trimmed, "#") {
            continue
        }

        parts := strings.Split(trimmed, ",")
        if len(parts) < 2 {
            errors = append(errors, fmt.Errorf("invalid format: %s", item))
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        if key == "" {
            errors = append(errors, fmt.Errorf("empty key for item: %s", item))
            continue
        }

        if value == "" {
            errors = append(errors, fmt.Errorf("empty value for key: %s", key))
            continue
        }

        validated, err := validateEntry(key, value)
        if err != nil {
            errors = append(errors, err)
            continue
        }

        transformed := transformEntry(validated)
        if transformed == "" {
            errors = append(errors, fmt.Errorf("transform failed for: %s", key))
            continue
        }

        normalized := normalizeEntry(transformed)
        if normalized == "" {
            errors = append(errors, fmt.Errorf("normalize failed for: %s", key))
            continue
        }

        deduplicated := deduplicateEntry(normalized, result)
        if deduplicated == "" {
            continue
        }

        result = append(result, deduplicated)
    }

    if len(errors) > 0 {
        return result, fmt.Errorf("encountered %d errors", len(errors))
    }

    return result, nil
}
```

### Valid

```golang
// This function is short and focused, well under the 60-line limit.
func processData(input []string) ([]string, error) {
    result := make([]string, 0, len(input))

    for _, item := range input {
        processed, err := processItem(item)
        if err != nil {
            return nil, err
        }
        if processed != "" {
            result = append(result, processed)
        }
    }

    return result, nil
}
```

```golang
// Helper functions keep each function short and focused.
func processItem(item string) (string, error) {
    trimmed := strings.TrimSpace(item)
    if trimmed == "" || strings.HasPrefix(trimmed, "#") {
        return "", nil
    }

    parts := strings.Split(trimmed, ",")
    if len(parts) < 2 {
        return "", fmt.Errorf("invalid format: %s", item)
    }

    return validateAndTransform(parts[0], parts[1])
}
```
