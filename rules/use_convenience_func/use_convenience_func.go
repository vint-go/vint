package use_convenience_func

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseConvenienceFuncRule detects function calls that can be replaced with
// convenience wrappers from the Go standard library.
type UseConvenienceFuncRule struct{}

// Apply applies the rule to given file.
func (r *UseConvenienceFuncRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintConvenienceFunc{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseConvenienceFuncRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintConvenienceFunc{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseConvenienceFuncRule) Name() string {
	return "useConvenienceFunc"
}

// Group returns the rule group.
func (*UseConvenienceFuncRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseConvenienceFuncRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// convenienceMapping describes a function call pattern where the last argument
// is a specific integer literal (-1) and the call can be replaced with a
// convenience wrapper.
type convenienceMapping struct {
	pkg           string
	funcName      string
	lastArgValue  string // the literal value (e.g. "-1")
	argCount      int    // expected total number of arguments
	replacement   string // the convenience function name (e.g. "strings.Split")
}

var convenienceMappings = []convenienceMapping{
	{pkg: "strings", funcName: "SplitN", lastArgValue: "-1", argCount: 3, replacement: "strings.Split"},
	{pkg: "strings", funcName: "SplitAfterN", lastArgValue: "-1", argCount: 3, replacement: "strings.SplitAfter"},
	{pkg: "strings", funcName: "Replace", lastArgValue: "-1", argCount: 4, replacement: "strings.ReplaceAll"},
	{pkg: "bytes", funcName: "SplitN", lastArgValue: "-1", argCount: 3, replacement: "bytes.Split"},
	{pkg: "bytes", funcName: "SplitAfterN", lastArgValue: "-1", argCount: 3, replacement: "bytes.SplitAfter"},
	{pkg: "bytes", funcName: "Replace", lastArgValue: "-1", argCount: 4, replacement: "bytes.ReplaceAll"},
}

type lintConvenienceFunc struct {
	onFailure func(lint.Failure)
}

func (w *lintConvenienceFunc) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, m := range convenienceMappings {
		if !astutils.IsPkgDotName(call.Fun, m.pkg, m.funcName) {
			continue
		}

		if len(call.Args) != m.argCount {
			continue
		}

		lastArg := call.Args[len(call.Args)-1]
		if !isNegativeOne(lastArg) {
			continue
		}

		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       call,
			Category:   lint.FailureCategoryStyle,
			Failure:    m.pkg + "." + m.funcName + " can be replaced with " + m.replacement,
		})

		return w
	}

	return w
}

// isNegativeOne checks if the expression is the integer literal -1.
// It handles the case where the AST represents -1 as a unary minus applied to 1.
func isNegativeOne(expr ast.Expr) bool {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	if unary.Op != token.SUB {
		return false
	}
	lit, ok := unary.X.(*ast.BasicLit)
	if !ok {
		return false
	}
	return lit.Kind == token.INT && lit.Value == "1"
}
