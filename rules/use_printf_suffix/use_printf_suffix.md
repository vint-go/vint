---
title: usePrintfSuffix
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/usePrintfSuffix`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/usePrintfSuffix:
    # no additional options
```

## Details

Checks that printf-like functions are named with `f` at the end.

In Go, it is a well-established convention that functions accepting a format string and variadic arguments (similar to `fmt.Printf`) should have names ending with the letter `f`. This convention makes it immediately clear to callers that the function performs formatting. Standard library examples include `fmt.Printf`, `fmt.Sprintf`, `fmt.Errorf`, `log.Printf`, and `log.Fatalf`.

This rule identifies functions that match the printf-like signature but do not follow this naming convention. A function is considered printf-like when it meets all of the following criteria:

1. It has no return values.
2. It has at least two parameters.
3. The second-to-last parameter is of type `string` and is named `format`.
4. The last parameter is variadic with type `...interface{}` or `...any`.

When such a function is found and its name does not end with `f`, the linter reports:

> printf-like formatting function 'name' should be named 'namef'

This naming convention helps tools like `go vet` and `staticcheck` automatically verify format string correctness at compile time, catching mismatched format verbs and arguments early.

Source: https://github.com/golangci/go-printf-func-name

## Examples

### Invalid

```golang
// Function has printf-like signature but name does not end with 'f'
func myLog(format string, args ...interface{}) {
    const prefix = "[my] "
    log.Printf(prefix+format, args...)
}
```

```golang
// Custom error logging function missing the 'f' suffix
func logError(format string, args ...any) {
    log.Printf("[ERROR] "+format, args...)
}
```

```golang
// Wrapper around fmt.Sprintf-style formatting without 'f' suffix
func customPrint(format string, args ...interface{}) {
    fmt.Printf("[custom] "+format, args...)
}
```

### Valid

```golang
// Function name correctly ends with 'f'
func myLogf(format string, args ...interface{}) {
    const prefix = "[my] "
    log.Printf(prefix+format, args...)
}
```

```golang
// Properly named custom error logging function
func logErrorf(format string, args ...any) {
    log.Printf("[ERROR] "+format, args...)
}
```

```golang
// Properly named wrapper function
func customPrintf(format string, args ...interface{}) {
    fmt.Printf("[custom] "+format, args...)
}
```

```golang
// Not a printf-like function (has return value) - not flagged
func formatMessage(format string, args ...interface{}) string {
    return fmt.Sprintf(format, args...)
}
```

```golang
// Not a printf-like function (parameter not named 'format') - not flagged
func logMessage(msg string, args ...interface{}) {
    log.Println(msg)
}
```

```golang
// Not a printf-like function (only one parameter) - not flagged
func printLine(msg string) {
    fmt.Println(msg)
}
```
