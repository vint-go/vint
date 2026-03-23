# Description of available rules

List of all available rules.

<!-- toc -->

- [Configuration options format](#configuration-options-format)
- [bare-return](#bare-return)
- [constant-logical-expr](#constant-logical-expr)
- [datarace](#datarace)
- [empty-block](#empty-block)
- [errorf](#errorf)
- [import-shadowing](#import-shadowing)
- [increment-decrement](#increment-decrement)
- [struct-tag](#struct-tag)
- [var-naming](#var-naming)

<!-- tocstop -->

## Configuration options format

By convention, rule configuration options are documented using the `kebab-case` format (e.g., `max-lit-count`, `allow-strs`, `skip-comments`).
For backward compatibility, `camelCase` (e.g., `maxLitCount`, `allowStrs`, `skipComments`)
and `lowercase` (e.g., `maxlitcount`, `allowstrs`, `skipcomments`) formats are still supported but are deprecated.

## bare-return

_Description_: Warns on bare (a.k.a. naked) returns

_Configuration_: N/A

## constant-logical-expr

_Description_: The rule spots logical expressions that evaluate always to the same value.

_Configuration_: N/A

## datarace

_Description_: This rule spots potential dataraces caused by goroutines capturing (by-reference) particular identifiers of the function from
which goroutines are created.
The rule is able to spot two of such cases: go-routines capturing named return values, and capturing `for-range` values.

_Configuration_: N/A


## empty-block

_Description_: Empty blocks make code less readable and could be a symptom of a bug or unfinished refactoring.

_Configuration_: N/A

### Limitations

The `empty-block` rule has limited support for detecting intentionally empty `for` loops where the loop body is empty but the
loop controls (`Init`, `Cond`, or `Post`) contain function calls that perform the actual work.

Currently, the rule only recognizes a narrow pattern:

```go
for process() {
    // Intentionally empty - process() does the work
}

for range processChan {
    // Intentionally empty - draining the channel
}
```

However, it will produce **false positives** for more complex patterns such as:

```go
// False positive: rule will warn even though this is intentional
for _, c := step(); c; _, c = step() {
    // Loop body is intentionally empty; step() does the work
}

// False positive: rule will warn even though this is intentional
for p := 0; bar(p); p++ {
    // Loop body is intentionally empty; bar(p) does the work
}
```

**Workaround**: If you have intentionally empty `for` loops with function calls in the loop controls, you can disable the rule
in-place using a directive comment:

```go
//revive:disable:empty-block
for _, c := step(); c; _, c = step() {
    // Intentionally empty - step() does the work
}
//revive:enable:empty-block
```

The reason for this limitation is that properly detecting whether a `for` loop is intentionally empty requires understanding the
**semantics** of the called functions (whether they have side effects, modify state, etc.), which is beyond the scope of static
analysis that this rule performs.

For more details, see:

- <https://github.com/vint-go/vint/issues/1622>
- <https://github.com/vint-go/vint/issues/386>

## import-shadowing

_Description_: In Go it is possible to declare identifiers (packages, structs,
interfaces, parameters, receivers, variables, constants...) that conflict with the
name of an imported package. This rule spots identifiers that shadow an import.

The rule ignores versioned import paths such as `k8s.io/api/core/v1` when `v1` is the package name,
which allows identifiers like `v1`. This is a deliberate trade-off to keep the rule simple.

_Configuration_: N/A

## increment-decrement

_Description_: By convention, for better readability, incrementing an integer variable by 1 is recommended to be done using the `++` operator.
This rule spots expressions like `i += 1` and `i -= 1` and proposes to change them into `i++` and `i--`.

_Configuration_: N/A

## struct-tag

_Description_: The rule spots errors in struct tags.
This is useful because struct tags are not checked at compile time.

The list of [supported tags](https://go.dev/wiki/Well-known-struct-tags):

| Tag           | Documentation                                                            |
| ------------- | ------------------------------------------------------------------------ |
| `asn1`         | <https://pkg.go.dev/encoding/asn1>                                      |
| `bson`         | <https://pkg.go.dev/go.mongodb.org/mongo-driver/bson>                   |
| `cbor`         | <https://pkg.go.dev/github.com/fxamacker/cbor/v2>                   |
| `datastore`    | <https://pkg.go.dev/cloud.google.com/go/datastore>                      |
| `default`      | The type of "default" must match the type of the field.                 |
| `json`         | <https://pkg.go.dev/encoding/json>                                      |
| `mapstructure` | <https://pkg.go.dev/github.com/mitchellh/mapstructure>                  |
| `properties`   | <https://pkg.go.dev/github.com/magiconair/properties#Properties.Decode> |
| `protobuf`     | <https://github.com/golang/protobuf>                                    |
| `required`     | Should be only "true" or "false".                                       |
| `spanner`      | <https://pkg.go.dev/cloud.google.com/go/spanner>                        |
| `toml`         | <https://pkg.go.dev/github.com/pelletier/go-toml/v2>                    |
| `url`          | <https://github.com/google/go-querystring>                              |
| `validate`     | <https://github.com/go-playground/validator>                            |
| `xml`          | <https://pkg.go.dev/encoding/xml>                                       |
| `yaml`         | <https://pkg.go.dev/gopkg.in/yaml.v2>                                   |

_Configuration_: (optional) The list of struct tags that can be accepted by the rule additionally to the supported tags.

Configuration example:

To accept the `inline` option in JSON tags (and `outline` and `gnu` in BSON tags) you must provide the following configuration

```toml
[rule.struct-tag]
arguments = ["json,inline", "bson,outline,gnu"]
```

To prevent a tag from being checked, simply add a `!` before its name.
For example, to instruct the rule not to check `validate` tags
(and accept `outline` and `gnu` in BSON tags) you can provide the following configuration

```toml
[rule.struct-tag]
arguments = ["!validate", "bson,outline,gnu"]
```

## var-naming

_Description_: This rule warns when [initialism](https://go.dev/wiki/CodeReviewComments#initialisms), [variable](https://go.dev/wiki/CodeReviewComments#variable-names)
naming conventions are not followed.
It ignores functions starting with `Example`, `Test`, `Benchmark`, and `Fuzz` in test files, preserving `golint` original behavior.

_Configuration_: This rule accepts two slices of strings and one optional slice containing a single map with named parameters.
(This is because TOML does not support "slice of any," and we maintain backward compatibility with the previous configuration version).
The first slice is an allowlist, and the second one is a blocklist of initialisms.
You can add a boolean parameter `skip-initialism-name-checks` to control how names
of functions, variables, consts, and structs handle known initialisms (e.g., JSON, HTTP, etc.) when written in `camelCase`.
When `skip-initialism-name-checks` is set to true, the rule allows names like `readJson`, `HttpMethod` etc.
In the map, you can add a boolean `upper-case-const` parameter to allow `UPPER_CASE` for `const`.

By default, the rule behaves exactly as the alternative in `golint` for non-package identifiers;
`golint`-equivalent package-name warnings now require enabling the [`package-naming`](#package-naming) rule.
The legacy package-related options `skip-package-name-checks`, `extra-bad-package-names`, and `skip-package-name-collision-with-go-std` are deprecated
and are now treated as no-ops by `var-naming` (they are ignored, apart from an optional warning when logging is enabled).
Package-name checks should be configured via the [`package-naming`](#package-naming) rule instead,
and these options should be removed from `var-naming` configurations to avoid confusion.

Configuration examples:

```toml
[rule.var-naming]
arguments = [[], [], [{ skip-initialism-name-checks = true }]]
```

```toml
[rule.var-naming]
arguments = [["ID"], ["VM"], [{ upper-case-const = true }]]
```

## noDuplicateCode

_Description_: Detects duplicate fragments of code across Go source files using suffix tree analysis on serialized ASTs.
Structurally identical code blocks are flagged regardless of specific variable names or literal values.

_Configuration_: (int) the minimum token threshold for duplicate detection (default `150`).

## noExcessiveStatements

_Description_: Checks that functions do not exceed a maximum number of statements (default: 40).
Counts executable statements recursively, including those inside control flow structures, inline function literals, `go` statements, and `defer` statements.

_Configuration_: (int) the maximum number of statements allowed per function (default `40`). Set to `-1` to disable.

## noHighCyclomaticComplexity

_Description_: Checks the cyclomatic complexity of Go functions and reports those that exceed a configurable threshold (default: 30).
Complexity increases by +1 for each `if`, `for`, `case`, `&&`, or `||`.

_Configuration_: (int) the maximum cyclomatic complexity allowed per function (default `30`).

## noLongFunctions

_Description_: Checks that functions do not exceed a maximum number of lines (default: 60).
Lines are counted from the opening brace to the closing brace of the function body, excluding the function signature.

_Configuration_: (int) the maximum number of lines allowed per function (default `60`). Set to `-1` to disable.

## noBlankErrorAssignment

_Description_: Detects when error return values are explicitly assigned to the blank identifier (`_`).
This rule is not enabled by default and must be explicitly enabled.

_Configuration_: N/A

## noDeniedImport

_Description_: Reports when a Go source file imports a package that appears on the deny list configured for the matching depguard rule group.
Package matching uses prefix matching by default; append `$` for exact matching.

_Configuration_: ([]map) list of rule group configurations with `deny`, `allow`, `files`, and `list-mode` fields.

## noDirectErrorComparison

_Description_: Flags direct comparisons of error values using `==` or `!=` and recommends using `errors.Is()` instead.
Comparisons to `nil` and `io.EOF` are allowed. An auto-fix is available.

_Configuration_: N/A

## noDynamicErrors

_Description_: Flags the creation of dynamic errors inside functions using `errors.New()` or `fmt.Errorf()` without `%w`, and requires errors be defined as package-level sentinel variables.
Wrapping with `fmt.Errorf` using `%w` is allowed.

_Configuration_: N/A

## noFileScopedDeniedImport

_Description_: Reports when a package import violates a file-scoped depguard rule, meaning the import is prohibited specifically in certain types of files based on glob patterns.
Commonly used to prevent test dependencies from leaking into production code.

_Configuration_: ([]map) list of rule group configurations with `files`, `deny`, and `allow` fields.

## noSpaceInDirective

_Description_: Detects Go compiler directives that contain a space between the comment slashes (`//`) and the `go:` prefix (e.g., `// go:embed` instead of `//go:embed`).
A space causes the compiler to silently ignore the directive.

_Configuration_: N/A
