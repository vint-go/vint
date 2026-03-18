package no_implicit_const_value

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoImplicitConstValueRule detects const groups where only the first constant
// has an explicit value and subsequent constants implicitly repeat it. This
// pattern is different from iota-based sequences and is often unintentional.
type NoImplicitConstValueRule struct{}

// Apply applies the rule to given file.
func (r *NoImplicitConstValueRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintImplicitConstValue{onFailure: onFailure}

	for _, decl := range file.AST.Decls {
		ast.Walk(w, decl)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImplicitConstValueRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintImplicitConstValue{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoImplicitConstValueRule) Name() string {
	return "noImplicitConstValue"
}

// Group returns the rule group.
func (*NoImplicitConstValueRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImplicitConstValueRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintImplicitConstValue struct {
	onFailure func(lint.Failure)
}

func (w *lintImplicitConstValue) Visit(node ast.Node) ast.Visitor {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return w
	}

	if genDecl.Tok != token.CONST {
		return w
	}

	// Only check parenthesized const groups with multiple specs
	if !genDecl.Lparen.IsValid() || len(genDecl.Specs) < 2 {
		return w
	}

	// Check the first spec: it must have explicit values
	firstSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
	if !ok || len(firstSpec.Values) == 0 {
		return w
	}

	// If the first spec uses iota anywhere, skip — iota-based patterns are intentional
	if containsIota(firstSpec) {
		return w
	}

	// Check subsequent specs for implicit values (no explicit values)
	for _, spec := range genDecl.Specs[1:] {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// If a spec has no values, it implicitly repeats the previous expression
		if len(vs.Values) == 0 {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 0.8,
				Node:       vs,
				Failure:    "constant implicitly repeats the previous value; consider making it explicit",
			})
		}
	}

	return w
}

// containsIota checks whether a ValueSpec contains a reference to iota.
func containsIota(spec *ast.ValueSpec) bool {
	for _, val := range spec.Values {
		if hasIota(val) {
			return true
		}
	}
	return false
}

// hasIota recursively checks whether an expression references iota.
func hasIota(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "iota"
	case *ast.BinaryExpr:
		return hasIota(e.X) || hasIota(e.Y)
	case *ast.UnaryExpr:
		return hasIota(e.X)
	case *ast.ParenExpr:
		return hasIota(e.X)
	case *ast.CallExpr:
		for _, arg := range e.Args {
			if hasIota(arg) {
				return true
			}
		}
	}
	return false
}
