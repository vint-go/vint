---
title: noHighCyclomaticComplexity
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/complexity/noHighCyclomaticComplexity`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/complexity/noHighCyclomaticComplexity:
    min-complexity: 30
```

## Details

Checks the cyclomatic complexity of Go functions and reports those that exceed a configurable threshold.

Cyclomatic complexity is a software metric that measures the number of linearly independent paths through a function's source code. It is computed as follows:

- Every function starts with a base complexity of **1**.
- The complexity increases by **+1** for each `if`, `for`, `case`, `&&`, or `||`.

A higher cyclomatic complexity value indicates that a function has more branching logic, making it harder to understand, test, and maintain. Functions with high cyclomatic complexity are strong candidates for refactoring into smaller, more focused functions.

The default threshold is **30**, meaning any function with a cyclomatic complexity greater than 30 will be reported. This threshold can be adjusted via the `min-complexity` configuration option.

Individual functions can be excluded from this check by placing a `//gocyclo:ignore` comment directive directly above the function declaration.

Source: https://github.com/fzipp/gocyclo

## Examples

### Invalid

```golang
// Cyclomatic complexity of 11 (exceeds threshold)
func processRequest(r *Request) error {
    if r == nil {
        return errors.New("nil request")
    }

    if r.Method == "GET" || r.Method == "HEAD" {
        if r.URL == nil {
            return errors.New("nil URL")
        }
        for _, header := range r.Headers {
            if header.Key == "Authorization" && header.Value != "" {
                if validateToken(header.Value) {
                    // process authorized request
                } else {
                    return errors.New("invalid token")
                }
            } else if header.Key == "Content-Type" {
                // handle content type
            }
        }
    } else if r.Method == "POST" || r.Method == "PUT" {
        if r.Body == nil {
            return errors.New("nil body")
        }
    }

    return nil
}
```

```golang
// Cyclomatic complexity of 12 (exceeds threshold)
func evaluate(op string, a, b int) (int, error) {
    switch op {
    case "add":
        return a + b, nil
    case "sub":
        return a - b, nil
    case "mul":
        return a * b, nil
    case "div":
        if b == 0 {
            return 0, errors.New("division by zero")
        }
        return a / b, nil
    case "mod":
        if b == 0 {
            return 0, errors.New("modulo by zero")
        }
        return a % b, nil
    case "pow":
        result := 1
        for i := 0; i < b; i++ {
            result *= a
        }
        return result, nil
    case "max":
        if a > b {
            return a, nil
        }
        return b, nil
    case "min":
        if a < b {
            return a, nil
        }
        return b, nil
    case "avg":
        return (a + b) / 2, nil
    default:
        return 0, fmt.Errorf("unknown operation: %s", op)
    }
}
```

### Valid

```golang
// Low cyclomatic complexity - clean, focused function
func add(a, b int) int {
    return a + b
}
```

```golang
// Moderate complexity within acceptable range
func greet(name string, formal bool) string {
    if formal {
        return "Hello, " + name + "."
    }
    return "Hi, " + name + "!"
}
```

```golang
// Function excluded from complexity check using ignore directive
//gocyclo:ignore
func complexButIgnored(input string) string {
    // This function will not be checked by gocyclo
    // regardless of its complexity
    if input == "" {
        return "empty"
    }
    for _, r := range input {
        if r == ' ' || r == '\t' {
            continue
        }
    }
    return input
}
```
