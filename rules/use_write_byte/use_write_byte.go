package use_write_byte

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseWriteByteRule detects WriteRune calls with a single-byte rune argument
// that could be replaced with WriteByte for better performance.
type UseWriteByteRule struct{}

// Apply applies the rule to the given file.
func (r *UseWriteByteRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintWriteByte{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseWriteByteRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintWriteByte{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseWriteByteRule) Name() string {
	return "useWriteByte"
}

// Group returns the rule group.
func (*UseWriteByteRule) Group() string {
	return "performance"
}

// CacheTier returns the cache tier for this rule.
func (*UseWriteByteRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintWriteByte struct {
	onFailure func(lint.Failure)
}

func (w *lintWriteByte) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Look for a method call: something.WriteRune(...)
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "WriteRune" {
		return w
	}

	// WriteRune should have exactly one argument
	if len(call.Args) != 1 {
		return w
	}

	// The argument should be a rune literal (CHAR token)
	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.CHAR {
		return w
	}

	// Check if the rune fits in a single byte (0-127 for ASCII)
	r, _, _, err := strconv.UnquoteChar(lit.Value[1:len(lit.Value)-1], '\'')
	if err != nil {
		return w
	}

	if r > 127 {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryOptimization,
		Failure:    "use WriteByte instead of WriteRune for single-byte rune argument",
	})

	return w
}
