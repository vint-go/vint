# goconst: ignore-calls setting not supported

## Affected rule
`lint/style/noRepeatedStrings` (maps to golangci-lint `goconst`)

## Behavior in golangci-lint
When goconst is configured with `ignore-calls: true` (the default), it ignores repeated string literals that appear only within function call arguments. This reduces noise by not flagging strings that are used as parameters (e.g. map keys, log messages) unless they also appear in other contexts like assignments or comparisons.

Example not flagged with `ignore-calls: true`:
```go
log.Println("some message")
log.Println("some message")
log.Println("some message")
```

## Behavior in vint
The vint `noRepeatedStrings` rule has no `ignore-calls` option. It reports all repeated string literals regardless of where they appear in the AST.

## Gap
The `ignore-calls` filtering is not available in vint. Users who relied on this golangci-lint default behavior may see additional findings after migration for strings that only appear in function call arguments.

## Example
```go
package example

import "fmt"

func example() {
    fmt.Println("repeated value")
    fmt.Println("repeated value")
    fmt.Println("repeated value")
}
```

golangci-lint with `ignore-calls: true` (default) does not flag this. vint's `noRepeatedStrings` will flag it.

## Impact on migration
- The `ignore-calls` setting is silently dropped during migration since there is no equivalent vint option.
- Users may see new false positives for repeated strings used only in function arguments.
- This is a low-severity gap since users can suppress individual findings or use the `ignore-strings` regex to exclude specific patterns.
