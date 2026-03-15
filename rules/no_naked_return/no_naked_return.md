---
title: noNakedReturn
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noNakedReturn`
- This rule is **recommended**, meaning it is enabled by default.
- This rule has a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noNakedReturn:
    maxFuncLines: 30
```

## Details

Disallow naked return statements in functions that exceed a configurable line length threshold.

A naked return is a `return` statement without explicit return values in a function that uses named return parameters. While naked returns are acceptable in short functions, they reduce readability and maintainability in longer functions because the reader must trace back to the function signature to understand what values are being returned.

As stated in Go's Code Review Comments: "Naked returns are okay if the function is a handful of lines. Once it's a medium sized function, be explicit with your return values."

This rule flags any naked return in a function whose body exceeds the configured maximum line count (default: 30 lines). The auto-fix replaces the naked `return` with an explicit return statement listing all named return parameters.

### Options

- `maxFuncLines` (default: `30`): The maximum number of lines a function may have while still being allowed to use naked returns. Functions exceeding this threshold will be flagged if they contain naked return statements.

Source: https://github.com/alexkohler/nakedret

## Examples

### Invalid

```golang
// A function with 35 lines and a naked return will be flagged.
func processData(input []byte) (result string, err error) {
    if len(input) == 0 {
        err = fmt.Errorf("empty input")
        return // naked return in a long function
    }

    decoded := string(input)

    if strings.HasPrefix(decoded, "error") {
        err = fmt.Errorf("invalid prefix")
        return // naked return in a long function
    }

    // ... many more lines of processing ...
    result = strings.TrimSpace(decoded)
    result = strings.ToLower(result)
    result = strings.ReplaceAll(result, "\n", " ")

    if len(result) > 1000 {
        result = result[:1000]
    }

    if strings.Contains(result, "forbidden") {
        err = fmt.Errorf("forbidden content")
        return // naked return in a long function
    }

    for i := 0; i < 10; i++ {
        result += fmt.Sprintf(" %d", i)
    }

    return // naked return in a long function
}
```

```golang
// A function literal (closure) exceeding the threshold with a naked return.
func main() {
    process := func() (count int, err error) {
        // ... many lines of code exceeding the threshold ...
        for i := 0; i < 100; i++ {
            count++
            if count > 50 {
                err = fmt.Errorf("too many")
                return // naked return in a long closure
            }
        }
        // ... additional lines ...
        count = count * 2
        count = count + 1
        count = count - 1
        count = count / 2
        count = count % 3
        count = count + 10
        count = count - 5
        count = count * 3
        count = count + 7
        count = count - 2
        count = count + 4
        count = count - 1
        count = count + 8
        count = count - 3
        count = count + 6
        count = count - 4
        count = count + 9
        count = count - 6
        count = count + 2
        count = count - 7
        return // naked return in a long closure
    }
    process()
}
```

### Valid

```golang
// Short function with a naked return is acceptable.
func add(a, b int) (sum int) {
    sum = a + b
    return
}
```

```golang
// Long function using explicit return values is acceptable.
func processData(input []byte) (result string, err error) {
    if len(input) == 0 {
        return "", fmt.Errorf("empty input")
    }

    decoded := string(input)

    if strings.HasPrefix(decoded, "error") {
        return "", fmt.Errorf("invalid prefix")
    }

    // ... many more lines of processing ...
    result = strings.TrimSpace(decoded)
    result = strings.ToLower(result)
    result = strings.ReplaceAll(result, "\n", " ")

    if len(result) > 1000 {
        result = result[:1000]
    }

    if strings.Contains(result, "forbidden") {
        return "", fmt.Errorf("forbidden content")
    }

    for i := 0; i < 10; i++ {
        result += fmt.Sprintf(" %d", i)
    }

    return result, nil
}
```

```golang
// Function without named return parameters cannot have naked returns.
func multiply(a, b int) int {
    return a * b
}
```
