package no_unnecessary_lambda

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnnecessaryLambdaRule detects function literals that can be simplified
// by replacing them with a direct reference to the wrapped function.
type NoUnnecessaryLambdaRule struct{}

// Apply applies the rule to given file.
func (r *NoUnnecessaryLambdaRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryLambda{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoUnnecessaryLambdaRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnnecessaryLambda{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoUnnecessaryLambdaRule) Name() string {
	return "noUnnecessaryLambda"
}

// Group returns the rule group.
func (*NoUnnecessaryLambdaRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnnecessaryLambdaRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnnecessaryLambda struct {
	onFailure func(lint.Failure)
}

func (w *lintUnnecessaryLambda) Visit(node ast.Node) ast.Visitor {
	funcLit, ok := node.(*ast.FuncLit)
	if !ok {
		return w
	}

	// The function literal body must contain exactly one statement
	if funcLit.Body == nil || len(funcLit.Body.List) != 1 {
		return w
	}

	// Extract the inner call from the single statement.
	// It can be either:
	// 1. An expression statement: func() { doWork() }
	// 2. A return statement with a single call: func() bool { return less(i, j) }
	var innerCall *ast.CallExpr

	switch stmt := funcLit.Body.List[0].(type) {
	case *ast.ExprStmt:
		call, ok := stmt.X.(*ast.CallExpr)
		if !ok {
			return w
		}
		// If the function literal has return values but body is just an expression
		// statement (not a return), this is not a simple wrapper.
		if funcLit.Type.Results != nil && len(funcLit.Type.Results.List) > 0 {
			return w
		}
		innerCall = call
	case *ast.ReturnStmt:
		// Must have exactly one return value that is a call
		if len(stmt.Results) != 1 {
			return w
		}
		call, ok := stmt.Results[0].(*ast.CallExpr)
		if !ok {
			return w
		}
		// The function literal must have return values for a return statement
		// to make sense as a simple wrapper
		if funcLit.Type.Results == nil || len(funcLit.Type.Results.List) == 0 {
			return w
		}
		innerCall = call
	default:
		return w
	}

	// The inner call must not use variadic expansion (...)
	if innerCall.Ellipsis.IsValid() {
		return w
	}

	// Get the function literal's parameters
	params := funcLit.Type.Params

	// Count total parameters
	paramCount := 0
	if params != nil {
		for _, field := range params.List {
			if len(field.Names) == 0 {
				paramCount++
			} else {
				paramCount += len(field.Names)
			}
		}
	}

	// The inner call must have the same number of arguments as the outer parameters
	if len(innerCall.Args) != paramCount {
		return w
	}

	// If there are parameters, each argument must be the corresponding parameter
	// passed in the same order
	if paramCount > 0 {
		paramNames := collectParamNames(params)
		for i, arg := range innerCall.Args {
			argIdent, ok := arg.(*ast.Ident)
			if !ok {
				return w
			}
			if argIdent.Name != paramNames[i] {
				return w
			}
		}
	}

	// Only flag when the inner call target is a plain identifier (package-level
	// function). Skip method calls (SelectorExpr), calls through variables, and
	// any other complex expressions — without type info we can't tell if the
	// callee is a stable function reference or a captured variable.
	funIdent, ok := innerCall.Fun.(*ast.Ident)
	if !ok {
		return w
	}

	// If the called identifier matches any of the lambda's own parameters,
	// it's a variable call (e.g. `func(c) error { return next(c) }` where
	// `next` is a captured parameter), not a direct function reference.
	if params != nil {
		for _, field := range params.List {
			for _, name := range field.Names {
				if name.Name == funIdent.Name {
					return w
				}
			}
		}
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       funcLit,
		Category:   lint.FailureCategoryStyle,
		Failure:    "unnecessary lambda, use the function directly",
	})

	return w
}

// collectParamNames returns the parameter names in order from a FieldList.
func collectParamNames(params *ast.FieldList) []string {
	var names []string
	for _, field := range params.List {
		if len(field.Names) == 0 {
			names = append(names, "_")
		} else {
			for _, name := range field.Names {
				names = append(names, name.Name)
			}
		}
	}
	return names
}
