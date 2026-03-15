---
title: noLineTooLong
---

## Summary

- Rule available since: `v0.0.1`
- Category: `lint/style/noLineTooLong`
- This rule is **recommended**, meaning it is enabled by default.
- This rule doesn't have a fix.
- The default severity of this rule is **warning**.

## Configuration
```yaml title="vint.yaml"
settings:
  lint/style/noLineTooLong:
    line-length: 120
    tab-width: 1
```

### Options

- `line-length` (default: `120`): The maximum allowed number of characters per line. Lines exceeding this limit will trigger a warning. The length is measured in Unicode runes (characters), not bytes.
- `tab-width` (default: `1`): The number of spaces that a tab character counts as when calculating line length. Each tab character in a line is replaced by this many spaces before measuring the line length.

## Details

Reports lines that exceed a configured maximum character length. Long lines reduce code readability and can make side-by-side diffs harder to review. This rule encourages developers to keep lines within a reasonable width.

The linter measures line length in Unicode rune count (not byte count), so multi-byte characters such as UTF-8 encoded symbols are counted as single characters.

The following lines are automatically excluded from the check:

- **Go compiler directives**: Lines starting with `//go:` (such as `//go:generate`, `//go:build`, `//go:noinline`, etc.) are skipped, since these directives can be long and are not easily breakable.
- **Multi-line import blocks**: Lines within `import (...)` blocks are skipped, because import paths are often long and cannot be meaningfully shortened.

Tab characters are expanded to spaces (using the configured `tab-width`) before measuring line length. This ensures consistent measurement regardless of editor tab settings.

Source: https://github.com/walle/lll

## Examples

### Invalid

```golang
// With default line-length of 120, this line is too long:
func processUserAccountDetailsAndReturnFormattedOutputWithExtendedValidationAndErrorHandlingForAllCases(ctx context.Context, userID string) (string, error) {
```

```golang
// A comment that is way too long and should be broken up into multiple lines to improve readability and maintain consistent code style across the project
```

```golang
var myMap = map[string]string{"key1": "value1", "key2": "value2", "key3": "value3", "key4": "value4", "key5": "value5", "key6": "value6"}
```

### Valid

```golang
// Short function signature that fits
// within the line length limit
func processUser(
	ctx context.Context,
	userID string,
) (string, error) {
	return "", nil
}
```

```golang
// Comment broken into multiple lines
// so each line stays within the
// configured maximum length.
```

```golang
// Go compiler directives are excluded
// (no matter how long):
//go:generate stringer -type=MyType -output=mytype_string.go -trimprefix=MyType
```

```golang
// Multi-line import blocks are excluded:
import (
	"github.com/some/very/long/package/path/that/would/normally/exceed/the/line/length/limit/but/is/excluded"
)
```
