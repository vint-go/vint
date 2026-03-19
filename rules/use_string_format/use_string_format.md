---
title: useStringFormat
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/useStringFormat`
- This rule is not recommended, meaning it is not enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/useStringFormat:
    arguments:
      - ["core.WriteError[1].Message", "/^([^A-Z]|$)/", "must not start with a capital letter"]
      - ["fmt.Errorf[0]", "/(^|[^\\.!?])$/", "must not end in punctuation"]
      - ["panic", "/^[^\\n]*$/", "must not contain line breaks"]
```

## Details

This rule allows you to configure a list of regular expressions that string literals in certain function calls are checked against. This is geared towards user-facing applications where string literals are often used for messages that will be presented to users, so it may be desirable to enforce consistent formatting.

Each argument is a slice containing 2-3 strings: a scope, a regex, and an optional error message.

1. The first string defines a **scope**. This controls which string literals the regex will apply to, and is defined as a function argument. It must contain at least a function name (`core.WriteError`). Scopes may optionally contain a number specifying which argument in the function to check (`core.WriteError[1]`), as well as a struct field (`core.WriteError[1].Message`, only works for top level fields). Function arguments are counted starting at 0. If no argument number is provided, the first argument will be used (same as `[0]`). You can use multiple scopes for one regex by splitting them with `,` (`core.WriteError,fmt.Errorf`).

2. The second string is a **regular expression** (beginning and ending with a `/` character), which will be used to check the string literals in the scope. The default semantics is "strings matching the regular expression are OK". If you need to inverse the semantics you can add a `!` just before the first `/`. For example, `"/^[A-Z].*$/"` accepts strings starting with capital letters, while `"!/^[A-Z].*$/"` fails on strings starting with capital letters.

3. The third string (optional) is a **message** containing the purpose for the regex, which will be used in lint errors.

Source: https://github.com/mgechev/revive

## Examples

### Invalid

```golang
// With config: ["fmt.Errorf[0]", "/^([^A-Z]|$)/", "must not start with a capital letter"]
func example() {
    fmt.Errorf("This starts with a capital letter") // lint error: must not start with a capital letter
}
```

```golang
// With config: ["panic", "/^[^\\n]*$/", "must not contain line breaks"]
func example() {
    panic("line one\nline two") // lint error: must not contain line breaks
}
```

### Valid

```golang
// With config: ["fmt.Errorf[0]", "/^([^A-Z]|$)/", "must not start with a capital letter"]
func example() {
    fmt.Errorf("this starts with a lowercase letter") // OK
}
```

```golang
// With config: ["panic", "/^[^\\n]*$/", "must not contain line breaks"]
func example() {
    panic("single line message") // OK
}
```
