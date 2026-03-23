package no_non_octal_file_mode

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoNonOctalFileModeRule detects non-octal integer literals used as file mode arguments,
// where the value looks like it was intended to be octal (e.g., 644 instead of 0644).
type NoNonOctalFileModeRule struct{}

// fileModeFunc describes a standard library function that takes an os.FileMode parameter.
type fileModeFunc struct {
	pkg      string
	name     string
	argIndex int // 0-based index of the FileMode argument
}

// knownFileModeFuncs lists standard library functions that accept os.FileMode.
var knownFileModeFuncs = []fileModeFunc{
	{"os", "OpenFile", 2},
	{"os", "Mkdir", 1},
	{"os", "MkdirAll", 1},
	{"os", "WriteFile", 2},
	{"os", "Chmod", 1},
}

// Apply applies the rule to given file.
func (r *NoNonOctalFileModeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNonOctalFileMode{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoNonOctalFileModeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintNonOctalFileMode{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoNonOctalFileModeRule) Name() string {
	return "noNonOctalFileMode"
}

// Group returns the rule group.
func (*NoNonOctalFileModeRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoNonOctalFileModeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNonOctalFileMode struct {
	onFailure func(lint.Failure)
}

func (w *lintNonOctalFileMode) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fmf := range knownFileModeFuncs {
		if !astutils.IsPkgDotName(ce.Fun, fmf.pkg, fmf.name) {
			continue
		}

		if fmf.argIndex >= len(ce.Args) {
			continue
		}

		arg := ce.Args[fmf.argIndex]
		lit, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}

		if lit.Kind != token.INT {
			continue
		}

		// Check if this is a decimal integer that looks like an octal file mode.
		// Octal literals start with "0" (e.g., 0644, 0o644, 0O644).
		// We want to flag decimal integers whose digits are all in [0-7],
		// suggesting the author intended octal notation.
		val := lit.Value

		// Skip if already octal (starts with 0), hex (0x), binary (0b), or is just "0"
		if strings.HasPrefix(val, "0") {
			continue
		}

		// Check if all digits are valid octal digits (0-7)
		if looksLikeOctalFileMode(val) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       lit,
				Category:   lint.FailureCategoryLogic,
				Failure:    fmt.Sprintf("file mode %s is not in octal; did you mean 0%s?", val, val),
			})
		}

		break
	}

	return w
}

// looksLikeOctalFileMode checks if a decimal integer string looks like it was
// intended to be an octal file mode (all digits are 0-7, non-empty, positive).
func looksLikeOctalFileMode(val string) bool {
	if len(val) == 0 {
		return false
	}

	for _, ch := range val {
		if ch < '0' || ch > '7' {
			return false
		}
	}

	return true
}
