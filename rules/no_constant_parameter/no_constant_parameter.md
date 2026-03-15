---
title: noConstantParameter
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/suspicious/noConstantParameter`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/suspicious/noConstantParameter:
    check-exported: false # Whether to check exported (public) functions. Default: false.
```

## Details

Reports function parameters that always receive the same constant value at every call site in the program. Using whole-program analysis, `unparam` inspects all callers of a function and determines if a parameter is invariantly passed the same argument. When a parameter always receives the same value, it can typically be replaced with a local constant or variable inside the function body, simplifying both the function signature and its call sites.

Parameters that always receive the same value are a form of unnecessary indirection. They add complexity to the API surface without providing any flexibility, since the value never varies. This often occurs after refactoring when a parameter that once varied has been narrowed down to a single usage pattern, but the signature was never updated to reflect that change.

By default, this rule only checks unexported (private) functions, since exported functions may intentionally accept a parameter for future extensibility or to satisfy an interface. Set `check-exported: true` to also analyze exported functions.

The typical diagnostic message is: `paramName always receives constValue`.

Source: https://github.com/mvdan/unparam

## Examples

### Invalid

```golang
// Every call to "repeat" passes 3 as the "times" argument.
func repeat(s string, times int) string {
    var result string
    for i := 0; i < times; i++ {
        result += s
    }
    return result
}

func main() {
    fmt.Println(repeat("ha", 3))
    fmt.Println(repeat("ho", 3))
    fmt.Println(repeat("he", 3))
}
```

```golang
// Every caller passes true for "verbose".
func logMessage(msg string, verbose bool) {
    if verbose {
        fmt.Println("[VERBOSE]", msg)
    } else {
        fmt.Println(msg)
    }
}

func run() {
    logMessage("starting", true)
    logMessage("processing", true)
    logMessage("done", true)
}
```

```golang
// The "separator" parameter is always ",".
func joinStrings(parts []string, separator string) string {
    return strings.Join(parts, separator)
}

func buildCSV(rows [][]string) string {
    var lines []string
    for _, row := range rows {
        lines = append(lines, joinStrings(row, ","))
    }
    return strings.Join(lines, "\n")
}
```

### Valid

```golang
// The "times" parameter receives different values at different call sites.
func repeat(s string, times int) string {
    var result string
    for i := 0; i < times; i++ {
        result += s
    }
    return result
}

func main() {
    fmt.Println(repeat("ha", 3))
    fmt.Println(repeat("ho", 5))
    fmt.Println(repeat("he", 1))
}
```

```golang
// The "verbose" parameter varies across call sites.
func logMessage(msg string, verbose bool) {
    if verbose {
        fmt.Println("[VERBOSE]", msg)
    } else {
        fmt.Println(msg)
    }
}

func run() {
    logMessage("starting", true)
    logMessage("processing", false)
    logMessage("done", true)
}
```

```golang
// The separator varies depending on the output format.
func joinStrings(parts []string, separator string) string {
    return strings.Join(parts, separator)
}

func buildOutput(rows [][]string, format string) string {
    var sep string
    if format == "csv" {
        sep = ","
    } else {
        sep = "\t"
    }
    var lines []string
    for _, row := range rows {
        lines = append(lines, joinStrings(row, sep))
    }
    return strings.Join(lines, "\n")
}
```
