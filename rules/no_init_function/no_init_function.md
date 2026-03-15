---
title: noInitFunction
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noInitFunction`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noInitFunction:
    # This rule has no additional configuration options.
```

## Details

Disallow the use of `init()` functions in Go packages.

In Go, `init()` functions are special functions that are automatically executed when a package is imported. While they can be convenient for performing setup tasks, they introduce implicit side effects that make code harder to understand, test, and maintain.

Problems with `init()` functions include:

- **Hidden side effects**: Importing a package silently triggers `init()` execution, which can lead to unexpected behavior and makes it difficult to reason about program startup order.
- **Testing difficulty**: Code that relies on `init()` functions is harder to test in isolation because the initialization logic runs automatically and cannot be easily controlled or mocked during tests.
- **Reduced readability**: Developers reading the code may not immediately notice that importing a package has side effects, leading to confusion about program behavior.
- **Tight coupling**: `init()` functions often create tight coupling between packages, making refactoring and dependency management more complex.
- **Non-deterministic ordering**: When multiple `init()` functions exist across packages, their execution order depends on import order, which can lead to subtle and hard-to-debug issues.

Instead of using `init()` functions, prefer explicit initialization by calling setup functions directly from `main()` or from a dedicated initialization path. This makes the program flow explicit and easier to follow.

Source: https://github.com/leighmcculloch/gochecknoinits

## Examples

### Invalid

```golang
package mypackage

func init() {
    // This will trigger the rule because init functions are not allowed.
    setupDatabase()
}
```

```golang
package mypackage

import "fmt"

func init() {
    fmt.Println("package initialized")
}

func init() {
    // Multiple init functions in the same file are also flagged.
    registerHandlers()
}
```

```golang
package config

var defaultConfig Config

func init() {
    defaultConfig = Config{
        Timeout: 30,
        Retries: 3,
    }
}
```

### Valid

```golang
package mypackage

// Use explicit initialization instead of init().
func Setup() {
    setupDatabase()
}
```

```golang
package mypackage

import "fmt"

// Call initialization explicitly from main or a startup function.
func Initialize() {
    fmt.Println("package initialized")
    registerHandlers()
}
```

```golang
package config

// Use a constructor function instead of init().
func NewDefaultConfig() Config {
    return Config{
        Timeout: 30,
        Retries: 3,
    }
}
```

```golang
package main

import "myapp/config"

func main() {
    // Explicit initialization makes the program flow clear.
    cfg := config.NewDefaultConfig()
    app := NewApp(cfg)
    app.Run()
}
```
