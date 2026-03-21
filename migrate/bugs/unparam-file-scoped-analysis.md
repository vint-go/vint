# unparam: file-scoped analysis vs whole-program SSA analysis

## Affected rules
- `lint/suspicious/noConstantParameter` (maps to golangci-lint `unparam`)
- `lint/suspicious/noConstantResult` (maps to golangci-lint `unparam`)
- `lint/suspicious/noUnusedParameter` (maps to golangci-lint `unparam`)

## Behavior in golangci-lint
The original `unparam` linter (github.com/mvdan/unparam) performs whole-program SSA (Static Single Assignment) analysis using call graph construction algorithms (CHA or RTA). It can:
- Analyze function parameters across the entire program, including cross-package calls
- Track parameter usage through the SSA intermediate representation
- Determine if a parameter is unused or always receives the same value across all call sites in the program
- Analyze methods (functions with receivers)
- Use configurable call graph algorithms (`algo: cha` for libraries, `algo: rta` for programs)

## Behavior in vint
Vint's rules use file-level AST analysis only:
- `noConstantParameter` inspects call sites within a single file and explicitly skips methods (functions with receivers)
- `noConstantResult` inspects return statements within a single function declaration
- `noUnusedParameter` inspects the function body for references to parameter identifiers

## Gap
1. **Cross-file/cross-package analysis**: If a function is called from multiple files or packages, vint will only see call sites in the same file. This can lead to both false positives (reporting a parameter as always-constant when other files pass different values) and false negatives (not detecting issues visible only through whole-program analysis).
2. **Methods**: `noConstantParameter` skips methods entirely, while the original unparam checks them.
3. **`algo` setting**: The golangci-lint `algo` setting (cha/rta) has no equivalent in vint and is silently ignored during migration.

## Example
```go
// file1.go
package mypackage

func process(mode string) { // unparam would see calls from file2.go too
    fmt.Println(mode)
}

func caller1() {
    process("fast")
}
```

```go
// file2.go
package mypackage

func caller2() {
    process("slow") // unparam sees this, vint's noConstantParameter does not
}
```

With the original unparam, `process` would NOT be flagged because `mode` receives different values ("fast" and "slow") across call sites. With vint's file-scoped analysis of file1.go alone, `mode` would be flagged as always receiving "fast".

## Impact on migration
- Users may see false positives from vint's rules that they did not see with unparam, particularly for functions called from multiple files.
- Users may see fewer true positives for cross-package parameter analysis.
- The `algo` setting in golangci-lint configuration is silently dropped during migration since it has no equivalent.
