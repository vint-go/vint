---
title: noUncheckedTypeAssertion
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/correctness/noUncheckedTypeAssertion`
- This rule is **not recommended**, meaning it is not enabled by default. You can enable it manually.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/correctness/noUncheckedTypeAssertion:
    enabled: true
```

## Details

Detects when type assertion results are not checked for success. In Go, a type assertion `x.(T)` can return a second boolean value indicating whether the assertion succeeded. If the assertion fails and the second value is not captured, the program will panic at runtime.

This check is disabled by default because some developers prefer to let type assertions panic when they are certain the assertion will succeed. However, enabling this rule enforces a safer coding practice where all type assertions use the comma-ok idiom to gracefully handle unexpected types.

A failed unchecked type assertion causes a runtime panic, which can crash the program or bring down a goroutine. Using the comma-ok idiom (`value, ok := x.(T)`) allows the program to handle the failure gracefully instead of panicking.

Source: https://github.com/kisielk/errcheck

## Examples

### Invalid

```golang
// Type assertion without checking the ok value -- will panic if i is not a string
package main

import "fmt"

func printString(i interface{}) {
    s := i.(string)
    fmt.Println(s)
}

func main() {
    printString("hello")
    printString(42) // This will panic at runtime
}
```

```golang
// Type assertion in a function argument without checking
package main

import "fmt"

func process(val interface{}) {
    fmt.Println(val.(int) + 1)
}
```

```golang
// Type assertion on an interface return value without checking
package main

import "fmt"

type Animal interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof" }

func getAnimal() Animal {
    return Dog{}
}

func main() {
    animal := getAnimal()
    dog := animal.(Dog)
    fmt.Println(dog.Speak())
}
```

### Valid

```golang
// Type assertion using the comma-ok idiom
package main

import "fmt"

func printString(i interface{}) {
    s, ok := i.(string)
    if !ok {
        fmt.Println("not a string")
        return
    }
    fmt.Println(s)
}

func main() {
    printString("hello")
    printString(42) // Handled gracefully
}
```

```golang
// Type assertion with error handling in function
package main

import "fmt"

func process(val interface{}) {
    num, ok := val.(int)
    if !ok {
        fmt.Println("expected an integer")
        return
    }
    fmt.Println(num + 1)
}
```

```golang
// Using a type switch instead of direct type assertion
package main

import "fmt"

func describe(i interface{}) {
    switch v := i.(type) {
    case string:
        fmt.Printf("String: %s\n", v)
    case int:
        fmt.Printf("Integer: %d\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

func main() {
    describe("hello")
    describe(42)
    describe(3.14)
}
```

```golang
// Type assertion with proper check on interface return
package main

import "fmt"

type Animal interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "Woof" }

func getAnimal() Animal {
    return Dog{}
}

func main() {
    animal := getAnimal()
    dog, ok := animal.(Dog)
    if !ok {
        fmt.Println("not a Dog")
        return
    }
    fmt.Println(dog.Speak())
}
```
