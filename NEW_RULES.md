# Implemented New Rules

## [lint/complexity/noDuplicateCode](lint-rules/no-duplicate-code.md)
Detects duplicate fragments of code across Go source files using suffix tree analysis on serialized ASTs. Structurally identical code blocks are flagged regardless of specific variable names or literal values.

## [lint/complexity/noExcessiveStatements](lint-rules/no-excessive-statements.md)
Checks that functions do not exceed a maximum number of statements (default: 40). Counts executable statements recursively, including those inside control flow structures, inline function literals, `go` statements, and `defer` statements.

## [lint/complexity/noHighCyclomaticComplexity](lint-rules/no-high-cyclomatic-complexity.md)
Checks the cyclomatic complexity of Go functions and reports those that exceed a configurable threshold (default: 30). Complexity increases by +1 for each `if`, `for`, `case`, `&&`, or `||`.

## [lint/complexity/noLongFunctions](lint-rules/no-long-functions.md)
Checks that functions do not exceed a maximum number of lines (default: 60). Lines are counted from the opening brace to the closing brace of the function body, excluding the function signature.

## [lint/correctness/noBlankErrorAssignment](lint-rules/no-blank-error-assignment.md)
Detects when error return values are explicitly assigned to the blank identifier (`_`). This rule is not enabled by default and must be explicitly enabled.

## [lint/correctness/noDeniedImport](lint-rules/no-denied-import.md)
Reports when a Go source file imports a package that appears on the deny list configured for the matching depguard rule group. Package matching uses prefix matching by default; append `$` for exact matching.

## [lint/correctness/noDirectErrorComparison](lint-rules/no-direct-error-comparison.md)
Flags direct comparisons of error values using `==` or `!=` and recommends using `errors.Is()` instead. Comparisons to `nil` and `io.EOF` are allowed.

## [lint/correctness/noDynamicErrors](lint-rules/no-dynamic-errors.md)
Flags the creation of dynamic errors inside functions using `errors.New()` or `fmt.Errorf()` without `%w`, and requires errors be defined as package-level sentinel variables. Wrapping with `fmt.Errorf` using `%w` is allowed.

## [lint/correctness/noFileScopedDeniedImport](lint-rules/no-file-scoped-denied-import.md)
Reports when a package import violates a file-scoped depguard rule, meaning the import is prohibited specifically in certain types of files based on glob patterns. Commonly used to prevent test dependencies from leaking into production code.

## [lint/correctness/noSpaceInDirective](lint-rules/no-space-in-directive.md)
Detects Go compiler directives that contain a space between the comment slashes (`//`) and the `go:` prefix (e.g., `// go:embed` instead of `//go:embed`). A space causes the compiler to silently ignore the directive.

## [lint/correctness/noUnallowedImport](lint-rules/no-unallowed-import.md)
Reports when a Go source file imports a package that is not present in the allow list configured for the matching depguard rule group. Particularly relevant in `strict` list mode, where every import must be explicitly allowed.

## [lint/correctness/noUncheckedError](lint-rules/no-unchecked-error.md)
Detects when error return values from function calls are silently ignored. A set of standard library functions whose error returns are commonly non-critical are excluded by default.

## [lint/correctness/noUncheckedTypeAssertion](lint-rules/no-unchecked-type-assertion.md)
Detects when type assertion results are not checked for success using the comma-ok idiom. A failed unchecked type assertion causes a runtime panic. This rule is not enabled by default.

## [lint/correctness/noRangeVariableAlias](lint-rules/no-range-variable-alias.md)
Detects implicit memory aliasing in `for...range` statements where taking the address of the loop variable creates pointers that all reference the same memory location. Applicable to Go 1.21 or lower.

