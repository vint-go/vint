# vint

Very fast linter for Golang. Designed to be MUCH FASTER replacement of golangci-lint and support as much rules as possible.

<p align="center">
  <img src="./assets/logo.png" alt="" width="300">
</p>

Here's how `vint` is different from `golangci-lint`:

- More than 7x faster running the same rules as golangci-lint.
- Clear centralized unified configuration.

This is still in progress, docs are not updated and in need of some cleanup, but vint already has close to 500 rules supported and benchmarks show quite well performance.

## Installation

TBD

## Usage

TBD

### Command Line Flags

TBD

### Sample Invocations

TBD

### Comment Directives

TBD

### Custom Configuration

TBD

### Rule-level file excludes

TBD

## 'legacy' Revive rules

List of all available rules. The rules ported from `golint` are left unchanged and indicated in the `golint` column.

| Name                  | Config | Description                                                      | `golint` | Typed |
| --------------------- | :----: | :--------------------------------------------------------------- | :------: | :---: |
| [`add-constant`](./RULES_DESCRIPTIONS.md#add-constant)        |  map   | Suggests using constant for magic numbers and string literals    |    no    |  no   |
| [`atomic`](./RULES_DESCRIPTIONS.md#atomic)              |  n/a   | Check for common mistaken usages of the `sync/atomic` package    |    no    |  no   |
| [`banned-characters`](./RULES_DESCRIPTIONS.md#banned-characters)          |  []string (defaults to []string{})   |  Checks banned characters in identifiers |    no    |  no   |
| [`bare-return`](./RULES_DESCRIPTIONS.md#bare-return) | n/a  | Warns on bare returns   |    no    |  no   |
| [`blank-imports`](./RULES_DESCRIPTIONS.md#blank-imports)       |  n/a   | Disallows blank imports                                          |   yes    |  no   |
| [`bool-literal-in-expr`](./RULES_DESCRIPTIONS.md#bool-literal-in-expr)|  n/a   | Suggests removing Boolean literals from logic expressions        |    no    |  no   |
| [`call-to-gc`](./RULES_DESCRIPTIONS.md#call-to-gc)   | n/a    | Warns on explicit call to the garbage collector    |    no    |  no   |
| [`cognitive-complexity`](./RULES_DESCRIPTIONS.md#cognitive-complexity)          |  int (defaults to 7) | Sets restriction for maximum Cognitive complexity.              |    no    |  no   |
| [`comment-spacings`](./RULES_DESCRIPTIONS.md#comment-spacings)          |  []string   |  Warns on malformed comments |    no    |  no   |
| [`comments-density`](./RULES_DESCRIPTIONS.md#comments-density) |  int (defaults to 0)  | Enforces a minimum comment / code relation |    no    |  no   |
| [`confusing-naming`](./RULES_DESCRIPTIONS.md#confusing-naming)    |  n/a   | Warns on methods with names that differ only by capitalization   |    no    |  no   |
| [`confusing-results`](./RULES_DESCRIPTIONS.md#confusing-results)   |  n/a   | Suggests to name potentially confusing function results          |    no    |  no   |
| [`constant-logical-expr`](./RULES_DESCRIPTIONS.md#constant-logical-expr)   |  n/a   | Warns on constant logical expressions                        |    no    |  no   |
| [`context-as-argument`](./RULES_DESCRIPTIONS.md#context-as-argument) |  n/a   | `context.Context` should be the first argument of a function.    |   yes    |  no   |
| [`cyclomatic`](./RULES_DESCRIPTIONS.md#cyclomatic)          |  int (defaults to 10)   | Sets restriction for maximum Cyclomatic complexity.              |    no    |  no   |
| [`datarace`](./RULES_DESCRIPTIONS.md#datarace)          |  n/a   |  Spots potential dataraces |    no    |  no   |
| [`deep-exit`](./RULES_DESCRIPTIONS.md#deep-exit)           |  n/a   | Looks for program exits in funcs other than `main()` or `init()` |    no    |  no   |
| [`defer`](./RULES_DESCRIPTIONS.md#defer)          |  map   |  Warns on some [defer gotchas](https://blog.learngoprogramming.com/5-gotchas-of-defer-in-go-golang-part-iii-36a1ab3d6ef1)       |    no    |  no   |
| [`dot-imports`](./RULES_DESCRIPTIONS.md#dot-imports)         |  n/a   | Forbids `.` imports.                                             |   yes    |  no   |
| [`duplicated-imports`](./RULES_DESCRIPTIONS.md#duplicated-imports) | n/a  | Looks for packages that are imported two or more times   |    no    |  no   |
| [`early-return`](./RULES_DESCRIPTIONS.md#early-return)          | []string   | Spots if-then-else statements where the predicate may be inverted to reduce nesting |    no    |  no   |
| [`empty-block`](./RULES_DESCRIPTIONS.md#empty-block)         |  n/a   | Warns on empty code blocks                                       |    no    |  no   |
| [`empty-lines`](./RULES_DESCRIPTIONS.md#empty-lines)   | n/a | Warns when there are heading or trailing newlines in a block              |    no    |  no   |
| [`epoch-naming`](./RULES_DESCRIPTIONS.md#epoch-naming)      |  n/a   | Enforces naming conventions for epoch time variables       |    no    |  yes   |
| [`enforce-map-style`](./RULES_DESCRIPTIONS.md#enforce-map-style) |  string (defaults to "any")  |  Enforces consistent usage of `make(map[type]type)` or `map[type]type{}` for map initialization. Does not affect `make(map[type]type, size)` constructions. |    no    |  no   |
| [`enforce-repeated-arg-type-style`](./RULES_DESCRIPTIONS.md#enforce-repeated-arg-type-style) |  string (defaults to "any")  |  Enforces consistent style for repeated argument and/or return value types. |    no    |  no   |
| [`enforce-slice-style`](./RULES_DESCRIPTIONS.md#enforce-slice-style) |  string (defaults to "any")  |  Enforces consistent usage of `make([]type, 0)` or `[]type{}` for slice initialization. Does not affect `make(map[type]type, non_zero_len, or_non_zero_cap)` constructions. |    no    |  no   |
| [`enforce-switch-style`](./RULES_DESCRIPTIONS.md#enforce-switch-style) |  []string (defaults to enforce occurrence and position) |  Enforces consistent usage of `default` on `switch` statements. |    no    |  no   |
| [`lint/style/useErrorNaming`](./rules/use_error_naming/use_error_naming.md)        |  n/a   | Naming of error variables.                                       |   yes    |  no   |
| [`error-return`](./RULES_DESCRIPTIONS.md#error-return)        |  n/a   | The error return parameter should be last.                       |   yes    |  no   |
| [`error-strings`](./RULES_DESCRIPTIONS.md#error-strings)       |  []string   | Conventions around error strings.                                |   yes    |  no   |
| [`errorf`](./RULES_DESCRIPTIONS.md#errorf)              |  n/a   | Should replace `errors.New(fmt.Sprintf())` with `fmt.Errorf()`   |   yes    |  yes  |
| [`exported`](./RULES_DESCRIPTIONS.md#exported)            |  []string   | Naming and commenting conventions on exported symbols.           |   yes    |  no   |
| [`file-header`](./RULES_DESCRIPTIONS.md#file-header)         | string (defaults to none)| Header which each file should have.                              |    no    |  no   |
| [`file-length-limit`](./RULES_DESCRIPTIONS.md#file-length-limit) | map (optional)| Enforces a maximum number of lines per file |    no    |  no   |
| [`filename-format`](./RULES_DESCRIPTIONS.md#filename-format) | regular expression (optional) | Enforces the formatting of filenames |   no    |  no   |
| [`flag-parameter`](./RULES_DESCRIPTIONS.md#flag-parameter)      |  n/a   | Warns on boolean parameters that create a control coupling       |    no    |  no   |
| [`forbidden-call-in-wg-go`](./RULES_DESCRIPTIONS.md#forbidden-call-in-wg-go)  |  n/a   | Warns on forbidden calls inside calls to wg.Go |    no    |  no   |
| [`function-length`](./RULES_DESCRIPTIONS.md#function-length)          |  int, int (defaults to 50 statements, 75 lines)   |  Warns on functions exceeding the statements or lines max |    no    |  no   |
| [`function-result-limit`](./RULES_DESCRIPTIONS.md#function-result-limit) |  int (defaults to 3)| Specifies the maximum number of results a function can return    |    no    |  no   |
| [`get-return`](./RULES_DESCRIPTIONS.md#get-return)          |  n/a   | Warns on getters that do not yield any result                    |    no    |  no   |
| [`identical-branches`](./RULES_DESCRIPTIONS.md#identical-branches)          |  n/a   | Spots if-then-else statements with identical `then` and `else` branches       |    no    |  no   |
| [`identical-ifelseif-branches`](./RULES_DESCRIPTIONS.md#identical-ifelseif-branches)          |  n/a   | Spots `if ... else if` chains with identical branches.       |    no    |  no   |
| [`identical-ifelseif-conditions`](./RULES_DESCRIPTIONS.md#identical-ifelseif-conditions)          |  n/a   | Spots identical conditions in  `if ... else if` chains.       |    no    |  no   |
| [`identical-switch-branches`](./RULES_DESCRIPTIONS.md#identical-switch-branches)          |  n/a   | Spots `switch` with identical branches.       |    no    |  no   |
| [`identical-switch-conditions`](./RULES_DESCRIPTIONS.md#identical-switch-conditions)          |  n/a   | Spots identical conditions in case clauses of `switch` statements.       |    no    |  no   |
| [`if-return`](./RULES_DESCRIPTIONS.md#if-return)           |  n/a   | Redundant if when returning an error.                            |   no    |  no   |
| [`import-alias-naming`](./RULES_DESCRIPTIONS.md#import-alias-naming)         | string or map[string]string (defaults to allow regex pattern `^[a-z][a-z0-9]{0,}$`) | Conventions around the naming of import aliases.                              |    no    |  no   |
| [`import-shadowing`](./RULES_DESCRIPTIONS.md#import-shadowing)   | n/a    | Spots identifiers that shadow an import    |    no    |  no   |
| [`imports-blocklist`](./RULES_DESCRIPTIONS.md#imports-blocklist)   | []string | Disallows importing the specified packages                     |    no    |  no   |
| [`increment-decrement`](./RULES_DESCRIPTIONS.md#increment-decrement) |  n/a   | Use `i++` and `i--` instead of `i += 1` and `i -= 1`.            |   yes    |  no   |
| [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)   |  []string   | Prevents redundant else statements.                              |   yes    |  no   |
| [`inefficient-map-lookup`](./RULES_DESCRIPTIONS.md#inefficient-map-lookup)   |  n/a  | Spots iterative searches for a key in a map           |   no    |  yes   |
| [`line-length-limit`](./RULES_DESCRIPTIONS.md#line-length-limit)   | int (defaults to 80) | Specifies the maximum number of characters in a line             |    no    |  no   |
| [`max-control-nesting`](./RULES_DESCRIPTIONS.md#max-control-nesting) |  int (defaults to 5)  | Sets restriction for maximum nesting of control structures. |    no    |  no   |
| [`max-public-structs`](./RULES_DESCRIPTIONS.md#max-public-structs)  |  int (defaults to 5)  | The maximum number of public structs in a file.                  |    no    |  no   |
| [`modifies-parameter`](./RULES_DESCRIPTIONS.md#modifies-parameter)  |  n/a   | Warns on assignments to function parameters                      |    no    |  no   |
| [`modifies-value-receiver`](./RULES_DESCRIPTIONS.md#modifies-value-receiver) |  n/a   | Warns on assignments to value-passed method receivers        |    no    |  yes  |
| [`nested-structs`](./RULES_DESCRIPTIONS.md#nested-structs)          |  n/a   |  Warns on structs within structs |    no    |  no   |
| [`optimize-operands-order`](./RULES_DESCRIPTIONS.md#optimize-operands-order)          |  n/a   |  Checks inefficient conditional expressions |    no    |  no   |
| [`usePackageComments`](./rules/use_package_comments/use_package_comments.md)    |  n/a   | Package commenting conventions.                                  |   yes    |  no   |
| [`package-naming`](./RULES_DESCRIPTIONS.md#package-naming)    |  map   | Checks that package names follow Go conventions and best practices |   no    |  no   |
| [`package-directory-mismatch`](./RULES_DESCRIPTIONS.md#package-directory-mismatch)    | string | Checks that package name matches containing directory name     |   no    |  no   |
| [`range-val-address`](./RULES_DESCRIPTIONS.md#range-val-address)|  n/a   | Warns if address of range value is used dangerously |    no    |  yes   |
| [`range-val-in-closure`](./RULES_DESCRIPTIONS.md#range-val-in-closure)|  n/a   | Warns if range value is used in a closure dispatched as goroutine|    no    |  no   |
| [`range`](./RULES_DESCRIPTIONS.md#range)               |  n/a   | Prevents redundant variables when iterating over a collection.   |   yes    |  no   |
| [`receiver-naming`](./RULES_DESCRIPTIONS.md#receiver-naming)     |  map (optional)   | Conventions around the naming of receivers.                      |   yes    |  no   |
| [`redefines-builtin-id`](./RULES_DESCRIPTIONS.md#redefines-builtin-id)|  n/a   | Warns on redefinitions of builtin identifiers                    |    no    |  no   |
| [`redundant-build-tag`](./RULES_DESCRIPTIONS.md#redundant-build-tag) | n/a   | Warns about redundant `// +build` comment lines |   no    |  no   |
| [`redundant-import-alias`](./RULES_DESCRIPTIONS.md#redundant-import-alias)          |  n/a   |  Warns on import aliases matching the imported package name |    no    |  no   |
| [`redundant-test-main-exit`](./RULES_DESCRIPTIONS.md#redundant-test-main-exit) |  n/a   | Suggests removing `Exit` call in `TestMain` function for test files |    no    |  no   |
| [`string-format`](./RULES_DESCRIPTIONS.md#string-format)          |  map   | Warns on specific string literals that fail one or more user-configured regular expressions            |    no    |  no   |
| [`string-of-int`](./RULES_DESCRIPTIONS.md#string-of-int)          |  n/a   | Warns on suspicious casts from int to string            |    no    |  yes   |
| [`struct-tag`](./RULES_DESCRIPTIONS.md#struct-tag)          |  []string   | Checks common struct tags like `json`, `xml`, `yaml`               |    no    |  no   |
| [`superfluous-else`](./RULES_DESCRIPTIONS.md#superfluous-else)    |  []string   | Prevents redundant else statements (extends [`indent-error-flow`](./RULES_DESCRIPTIONS.md#indent-error-flow)) |    no    |  no   |
| [`time-date`](./RULES_DESCRIPTIONS.md#time-date)         |  n/a   | Reports bad usage of `time.Date`.                 |   no    |  no   |
| [`time-equal`](./RULES_DESCRIPTIONS.md#time-equal)         |  n/a   | Suggests to use `time.Time.Equal` instead of `==` and `!=` for equality check time.                 |   no    |  yes  |
| [`time-naming`](./RULES_DESCRIPTIONS.md#time-naming)         |  n/a   | Conventions around the naming of time variables.                 |   yes    |  yes  |
| [`unchecked-type-assertion`](./RULES_DESCRIPTIONS.md#unchecked-type-assertion)         |  n/a   | Disallows type assertions without checking the result.                 |   no    |  no   |
| [`unconditional-recursion`](./RULES_DESCRIPTIONS.md#unconditional-recursion)          |  n/a   | Warns on function calls that will lead to (direct) infinite recursion |    no    |  no   |
| [`unexported-naming`](./RULES_DESCRIPTIONS.md#unexported-naming)          |  n/a   |  Warns on wrongly named un-exported symbols       |    no    |  no   |
| [`unexported-return`](./RULES_DESCRIPTIONS.md#unexported-return)   |  n/a   | Warns when a public return is from unexported type.              |   yes    |  yes  |
| [`unhandled-error`](./RULES_DESCRIPTIONS.md#unhandled-error)   | []string   | Warns on unhandled errors returned by function calls    |    no    |  yes   |
| [`unnecessary-if`](./RULES_DESCRIPTIONS.md#unnecessary-if)    |  n/a   | Identifies `if-else` statements that can be replaced by simpler statements  |    no    |  no   |
| [`unnecessary-format`](./RULES_DESCRIPTIONS.md#unnecessary-format)    |  n/a   | Identifies calls to formatting functions where the format string does not contain any formatting verbs          |    no    |  no   |
| [`unnecessary-stmt`](./RULES_DESCRIPTIONS.md#unnecessary-stmt)    |  n/a   | Suggests removing or simplifying unnecessary statements          |    no    |  no   |
| [`unreachable-code`](./RULES_DESCRIPTIONS.md#unreachable-code)    |  n/a   | Warns on unreachable code                                        |    no    |  no   |
| [`unsecure-url-scheme`](./RULES_DESCRIPTIONS.md#unsecure-url-scheme)          |  n/a   |  Checks for usage of potentially unsecure URL schemes. |    no    |  no   |
| [`unused-parameter`](./RULES_DESCRIPTIONS.md#unused-parameter)    |  n/a   | Suggests to rename or remove unused function parameters          |    no    |  no   |
| [`unused-receiver`](./RULES_DESCRIPTIONS.md#unused-receiver)   | n/a    | Suggests to rename or remove unused method receivers    |    no    |  no   |
| [`use-any`](./RULES_DESCRIPTIONS.md#use-any)          |  n/a   |  Proposes to replace `interface{}` with its alias `any` |    no    |  no   |
| [`use-errors-new`](./RULES_DESCRIPTIONS.md#use-errors-new) | n/a   | Spots calls to `fmt.Errorf` that can be replaced by `errors.New` |   no    |  no   |
| [`useFmtPrint`](./rules/use_fmt_print/use_fmt_print.md) | n/a   | Proposes to replace calls to built-in `print` and `println` with their equivalents from `fmt`. |   no    |  no   |
| [`use-slices-sort`](./RULES_DESCRIPTIONS.md#use-slices-sort) | n/a   | Proposes to replace calls to `sort.Ints`, `sort.Strings` and the like with their equivalents from `slices` package. |   no    |  no   |
| [`use-waitgroup-go`](./RULES_DESCRIPTIONS.md#use-waitgroup-go)          |  n/a   |  Proposes to replace `wg.Add ... go {... wg.Done ...}` idiom with `wg.Go` |    no    |  no   |
| [`useless-break`](./RULES_DESCRIPTIONS.md#useless-break)          |  n/a   |  Warns on useless `break` statements in case clauses |    no    |  no   |
| [`useless-fallthrough`](./RULES_DESCRIPTIONS.md#useless-fallthrough)  |  n/a   |  Warns on useless `fallthrough` statements in case clauses |    no    |  no   |
| [`var-declaration`](./RULES_DESCRIPTIONS.md#var-declaration)     |  n/a   | Reduces redundancies around variable declaration.                |   yes    |  yes  |
| [`var-naming`](./RULES_DESCRIPTIONS.md#var-naming)          |  allowlist & blocklist of initialisms   | Naming rules.                                                    |   yes    |  no   |
| [`waitgroup-by-value`](./RULES_DESCRIPTIONS.md#waitgroup-by-value)  |  n/a   | Warns on functions taking sync.WaitGroup as a by-value parameter |    no    |  no   |

## Configurable rules

Here you can find how you can configure some existing rules.
Rule configuration options are documented using the `kebab-case` format (e.g., `max-lit-count`, `allow-strs`, `skip-comments`), and this format is preferred.
For backward compatibility, `camelCase` (e.g., `maxLitCount`, `allowStrs`, `skipComments`)
and `lowercase` (e.g., `maxlitcount`, `allowstrs`, `skipcomments`) formats are still supported but are deprecated.

### `var-naming`

This rule accepts two slices of strings, an allowlist and a blocklist of initialisms. By default, the rule behaves exactly as the alternative
in `golint` but optionally, you can relax it (see [golint/lint/issues/89](https://github.com/golang/lint/issues/89))

```toml
[rule.var-naming]
arguments = [["ID"], ["VM"]]
```

This way, revive will not warn for an identifier called `customId` but will warn that `customVm` should be called `customVM`.

## Available Formatters

This section lists all the available formatters and provides a screenshot for each one.

### Friendly

![Friendly formatter](/assets/formatter-friendly.png)

### Stylish

![Stylish formatter](/assets/formatter-stylish.png)

### Default

The default formatter produces the same output as `golint`.

![Default formatter](/assets/formatter-default.png)

### Plain

The plain formatter produces the same output as the default formatter and appends the URL to the rule description.

![Plain formatter](/assets/formatter-plain.png)

### Unix

The unix formatter produces the same output as the default formatter but surrounds the rules in `[]`.

![Unix formatter](/assets/formatter-unix.png)

### JSON

The `json` formatter produces output in JSON format.

### NDJSON

The `ndjson` formatter produces output in [`Newline Delimited JSON`](https://github.com/ndjson/ndjson-spec) format.

### Checkstyle

The `checkstyle` formatter produces output in a [Checkstyle-like](https://checkstyle.sourceforge.io/) format.

### SARIF

The `sarif`  formatter produces output in SARIF, for _Static Analysis Results Interchange Format_,
a standard JSON-based format for the output of static analysis tools defined and promoted by [OASIS](https://www.oasis-open.org/).

Current supported version of the standard is [SARIF-v2.1.0](https://docs.oasis-open.org/sarif/sarif/v2.1.0/csprd01/sarif-v2.1.0-csprd01.html
).

## Extensibility

The tool can be extended with custom rules or formatters. This section contains additional information on how to implement such.

To extend the linter with a custom rule you can push it to this repository or use `revive` as a library (see below)

To add a custom formatter you'll have to push it to this repository or fork it.
This is due to the limited `-buildmode=plugin` support which [works only on Linux (with known issues)](https://pkg.go.dev/plugin).

### Writing a Custom Rule

See [DEVELOPING.md](./DEVELOPING.md) for instructions on how to write a custom rule.

#### Using `revive` as a library

If a rule is specific to your use case
(i.e. it is not a good candidate to be added to `revive`'s rule set) you can add it to your linter using `revive` as a linting engine.

The following code shows how to use `revive` in your application.
In the example only one rule is added (`myRule`), of course, you can add as many as you need to.
Your rules can be configured programmatically or with the standard `revive` configuration file.
The full rule set of `revive` is also actionable by your application.

```go
package main

import (
	"github.com/vint-go/vint/cli"
	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/revivelib"
)

func main() {
	cli.RunRevive(revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}))
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	// ...
}
```

You can still go further and use `revive` without its CLI, as part of your library, or your CLI:

```go
package mylib

import (
	"github.com/vint-go/vint/config"
	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/revivelib"
)

// Error checking removed for clarity
func LintMyFile(file string) {
	conf, _ := config.GetConfig("../defaults.toml")

	revive, _ := revivelib.New(
		conf, // Configuration file
		true, // Set exit status
		2048, // Max open files

		// Then add as many extra rules as you need
		revivelib.NewExtraRule(&myRule{}, lint.RuleConfig{}),
	)

	failuresChan, err := revive.Lint(
		revivelib.Include(file),
		revivelib.Exclude("./fixtures"),
		// You can use as many revivelib.Include or revivelib.Exclude as required
	)
	if err != nil {
		panic("Shouldn't have failed: " + err.Error())
	}

	// Now let's return the formatted errors
	failures, exitCode, _ := revive.Format("stylish", failuresChan)

	// failures is the string with all formatted lint error messages
	// exit code is 0 if no errors, 1 if errors (unless config options change it)
	// ... do something with them
}

type myRule struct{}

func (f myRule) Name() string {
	return "myRule"
}

func (f myRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	// ...
}
```

### Custom Formatter

Each formatter needs to implement the following interface:

```go
type Formatter interface {
	Format(<-chan Failure, Config) (string, error)
	Name() string
}
```

The `Format` method accepts a channel of `Failure` instances and the configuration of the enabled rules.
The `Name()` method should return a string different from the names of the already existing rules.
This string is used when specifying the formatter when invoking the `revive` CLI tool.

For a sample formatter, take a look at [this file](/formatter/json.go).

## Speed Comparison

Compared to `golint`, `revive` performs better because it lints the files for each individual rule into a separate goroutine.
Here's a basic performance benchmark on MacBook Pro Early 2013 run on Kubernetes:

### golint

```shellsession
$ time golint kubernetes/... > /dev/null

real    0m54.837s
user    0m57.844s
sys     0m9.146s
```

### revive's speed

```shellsession
# no type checking
$ time revive -config untyped.toml kubernetes/... > /dev/null

real    0m8.471s
user    0m40.721s
sys     0m3.262s
```

Keep in mind that if you use rules that require type checking, the performance may drop to 2x faster than `golint`:

```shellsession
# type checking enabled
$ time revive kubernetes/... > /dev/null

real    0m26.211s
user    2m6.708s
sys     0m17.192s
```

Currently, type-checking is enabled by default. If you want to run the linter without type-checking, remove all typed rules from the configuration file.

## Overriding colorization detection

By default, `revive` determines whether or not to colorize its output based on whether it's connected to a TTY or not.
This works for most use cases, but may not behave as expected if you use `revive` in a pipeline of commands,
where STDOUT is being piped to another command.

To force colorization, add `REVIVE_FORCE_COLOR=1` to the environment you're running in. For example:

```shell
REVIVE_FORCE_COLOR=1 revive -formatter friendly ./... | tee revive.log
```

_Open a PR to add your project_.

## Contributors

### Credits

This project is a fork of [revive](https://github.com/mgechev/revive) and contains all its original code and rules, with the addition of rules reproducing the behavior of these linters:

| | | | | |
| --- | --- | --- | --- | --- |
| [asasalint](https://github.com/alingse/asasalint) | [asciicheck](https://github.com/golangci/asciicheck) | [bidichk](https://github.com/breml/bidichk) | [canonicalheader](https://github.com/lasiar/canonicalheader) | [containedctx](https://github.com/sivchari/containedctx) |
| [contextcheck](https://github.com/kkHAIKE/contextcheck) | [cyclop](https://github.com/bkielbasa/cyclop) | [decorder](https://gitlab.com/bosi/decorder) | [dupword](https://github.com/Abirdcfly/dupword) | [durationcheck](https://github.com/charithe/durationcheck) |
| [errchkjson](https://github.com/breml/errchkjson) | [errname](https://github.com/Antonboom/errname) | [exhaustive](https://github.com/nishanths/exhaustive) | [exptostd](https://github.com/ldez/exptostd) | [fatcontext](https://github.com/Crocmagnon/fatcontext) |
| [forbidigo](https://github.com/ashanbrown/forbidigo) | [forcetypeassert](https://github.com/gostaticanalysis/forcetypeassert) | [gci](https://github.com/daixiang0/gci) | [ginkgolinter](https://github.com/nunnatsa/ginkgolinter) | [gochecknoglobals](https://github.com/leighmcculloch/gochecknoglobals) |
| [gochecksumtype](https://github.com/alecthomas/go-check-sumtype) | [gocognit](https://github.com/uudashr/gocognit) | [godot](https://github.com/tetafro/godot) | [gofmt](https://github.com/golangci/gofmt) | [gofumpt](https://github.com/mvdan/gofumpt) |
| [goheader](https://github.com/denis-tingaikin/go-header) | [goimports](https://github.com/golang/tools) | [golint](https://github.com/golang/lint) | [gomoddirectives](https://github.com/ldez/gomoddirectives) | [gomodguard](https://github.com/ryancurrah/gomodguard) |
| [gosimple](https://github.com/dominikh/go-tools) | [gosmopolitan](https://github.com/xen0n/gosmopolitan) | [grouper](https://github.com/leonklingele/grouper) | [iface](https://github.com/uudashr/iface) | [importas](https://github.com/julz/importas) |
| [inamedparam](https://github.com/macabu/inamedparam) | [interfacebloat](https://github.com/sashamelentyev/interfacebloat) | [kubeapilinter](https://github.com/kubernetes-sigs/kube-api-linter) | [logcheck](https://github.com/timonwong/loggercheck) | [makezero](https://github.com/ashanbrown/makezero) |
| [mirror](https://github.com/butuzov/mirror) | [modernize](https://github.com/golang/tools) | [musttag](https://github.com/go-simpler/musttag) | [nestif](https://github.com/nakabonne/nestif) | [nilerr](https://github.com/gostaticanalysis/nilerr) |
| [nilnesserr](https://github.com/alingse/nilnesserr) | [nilnil](https://github.com/Antonboom/nilnil) | [nonamedreturns](https://github.com/firefart/nonamedreturns) | [nosprintfhostport](https://github.com/stbenjam/no-sprintf-host-port) | [paralleltest](https://github.com/kunwardeep/paralleltest) |
| [perfsprint](https://github.com/catenacyber/perfsprint) | [prealloc](https://github.com/alexkohler/prealloc) | [predeclared](https://github.com/nishanths/predeclared) | [promlinter](https://github.com/yeya24/promlinter) | [protogetter](https://github.com/ghostiam/protogetter) |
| [reassign](https://github.com/curioswitch/go-reassign) | [recvcheck](https://github.com/raeperd/recvcheck) | [rowserrcheck](https://github.com/jingyugao/rowserrcheck) | [scopelint](https://github.com/kyoh86/scopelint) | [sorted](https://github.com/ravsii/sorted) |
| [spancheck](https://github.com/jjti/go-spancheck) | [sqlclosecheck](https://github.com/ryanrolds/sqlclosecheck) | [structcheck](https://github.com/golangci/check) | [stylecheck](https://github.com/dominikh/go-tools) | [tagliatelle](https://github.com/ldez/tagliatelle) |
| [testableexamples](https://github.com/maratori/testableexamples) | [thelper](https://github.com/kulti/thelper) | [tparallel](https://github.com/moricho/tparallel) | [typecheck](https://github.com/golangci/golangci-lint) | [usestdlibvars](https://github.com/sashamelentyev/usestdlibvars) |
| [usetesting](https://github.com/ldez/usetesting) | [varcheck](https://github.com/golangci/check) | [wastedassign](https://github.com/sanposhiho/wastedassign) | [wrapcheck](https://github.com/tomarrell/wrapcheck) | [wsl_v5](https://github.com/bombsimon/wsl) |
| [zerologlint](https://github.com/ykadowak/zerologlint) | | | | |


## License

MIT
