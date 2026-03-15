package no_swapped_arguments

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSwappedArgumentsRule detects function calls where arguments appear to be
// in the wrong order, such as passing a string literal as the first argument
// to strings.HasPrefix instead of the string to search.
type NoSwappedArgumentsRule struct{}

// Apply applies the rule to given file.
func (r *NoSwappedArgumentsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSwappedArguments{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoSwappedArgumentsRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSwappedArguments{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoSwappedArgumentsRule) Name() string {
	return "noSwappedArguments"
}

// Group returns the rule group.
func (*NoSwappedArgumentsRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoSwappedArgumentsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// suspiciousFunc describes a function where passing a literal as the first
// argument and a non-literal as the second is suspicious (arguments likely swapped).
type suspiciousFunc struct {
	pkg  string
	name string
	// minArgs is the minimum number of arguments the function takes
	// (we only check the first two).
	minArgs int
}

// suspiciousFuncs lists standard library functions where the first argument
// is typically a variable (the haystack) and the second is a literal
// (the needle/pattern). If a string literal appears as the first argument
// while the second is a non-literal, the arguments are likely swapped.
var suspiciousFuncs = []suspiciousFunc{
	// strings package
	{"strings", "HasPrefix", 2},
	{"strings", "HasSuffix", 2},
	{"strings", "Contains", 2},
	{"strings", "EqualFold", 2},
	{"strings", "Split", 2},
	{"strings", "SplitAfter", 2},
	{"strings", "SplitN", 3},
	{"strings", "SplitAfterN", 3},
	{"strings", "Count", 2},
	{"strings", "Index", 2},
	{"strings", "Replace", 4},
	{"strings", "ReplaceAll", 3},
	{"strings", "TrimPrefix", 2},
	{"strings", "TrimSuffix", 2},
	{"strings", "Cut", 2},
	// bytes package
	{"bytes", "HasPrefix", 2},
	{"bytes", "HasSuffix", 2},
	{"bytes", "Contains", 2},
	{"bytes", "EqualFold", 2},
	{"bytes", "Split", 2},
	{"bytes", "SplitAfter", 2},
	{"bytes", "SplitN", 3},
	{"bytes", "SplitAfterN", 3},
	{"bytes", "Count", 2},
	{"bytes", "Index", 2},
	{"bytes", "Replace", 4},
	{"bytes", "ReplaceAll", 3},
	{"bytes", "TrimPrefix", 2},
	{"bytes", "TrimSuffix", 2},
	{"bytes", "Cut", 2},
}

type lintSwappedArguments struct {
	onFailure func(lint.Failure)
}

func (w *lintSwappedArguments) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	for _, sf := range suspiciousFuncs {
		if pkgIdent.Name != sf.pkg || sel.Sel.Name != sf.name {
			continue
		}

		if len(ce.Args) < sf.minArgs {
			continue
		}

		// Check if first argument is a string literal and second is not
		if isStringLiteral(ce.Args[0]) && !isStringLiteral(ce.Args[1]) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryLogic,
				Failure:    sf.pkg + "." + sf.name + " arguments order looks suspicious, the literal string argument should probably not be first",
			})
		}

		break
	}

	return w
}

// isStringLiteral returns true if the expression is a string literal.
func isStringLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	return ok && lit.Kind == token.STRING
}
