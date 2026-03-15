---
title: noRepeatedNumbers
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noRepeatedNumbers`
- This rule is **not recommended**, meaning it is not enabled by default. It must be explicitly enabled in your configuration.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noRepeatedNumbers:
    enabled: true
    min-occurrences: 3
    min: 0
    max: 0
```

## Details

Detects numeric literals (integers and floats) that appear multiple times in the codebase and could be replaced by named constants.

Magic numbers scattered throughout the code reduce readability and maintainability. When a numeric value appears in multiple places, its meaning is often unclear, and changing it requires finding and updating every occurrence. Extracting repeated numbers into named constants makes the code self-documenting and easier to maintain.

This check is enabled via the `-numbers` flag in goconst or `numbers: true` in golangci-lint configuration. You can restrict the range of numbers to check using the `min` and `max` options, and control the reporting threshold with `min-occurrences`.

Common numbers like 0 and 1 are often excluded from this check because they are so ubiquitous that requiring constants for them would be counterproductive. Use the `min` option to set a lower bound (e.g., `min: 2`) to avoid reporting trivially common numbers.

Source: https://github.com/jgautheron/goconst

## Examples

### Invalid

```golang
// The number 3600 appears multiple times and should be a constant.
func CacheExpiry() time.Duration {
    return time.Duration(3600) * time.Second
}

func SessionExpiry() time.Duration {
    return time.Duration(3600) * time.Second
}

func TokenExpiry() time.Duration {
    return time.Duration(3600) * time.Second
}
```

```golang
// The number 100 appears multiple times.
func CalculatePercentage(value, total float64) float64 {
    return (value / total) * 100
}

func ScaleValue(value float64) float64 {
    return value * 100
}

func FormatPercentage(value float64) string {
    return fmt.Sprintf("%.2f%%", value*100)
}
```

### Valid

```golang
// The repeated number is extracted into a named constant.
const secondsPerHour = 3600

func CacheExpiry() time.Duration {
    return time.Duration(secondsPerHour) * time.Second
}

func SessionExpiry() time.Duration {
    return time.Duration(secondsPerHour) * time.Second
}

func TokenExpiry() time.Duration {
    return time.Duration(secondsPerHour) * time.Second
}
```

```golang
// The repeated number is extracted into a named constant.
const percentageMultiplier = 100

func CalculatePercentage(value, total float64) float64 {
    return (value / total) * percentageMultiplier
}

func ScaleValue(value float64) float64 {
    return value * percentageMultiplier
}
```

```golang
// Numbers that appear only once are fine.
func GetBufferSize() int {
    return 4096
}
```
