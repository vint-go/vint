package no_malformed_test_function

import (
	"fmt"
	"go/ast"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoMalformedTestFunctionRule checks that test, benchmark, fuzz, and example
// functions follow the correct naming conventions and signatures required by
// the go test command.
type NoMalformedTestFunctionRule struct{}

// Apply applies the rule to the given file.
func (r *NoMalformedTestFunctionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect all top-level exported identifiers for example function validation.
	topLevelIdents := collectTopLevelIdents(file.AST)

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Skip methods (functions with receivers).
		if funcDecl.Recv != nil {
			continue
		}

		name := funcDecl.Name.Name

		switch {
		case strings.HasPrefix(name, "Test"):
			if f := checkTestFunc(funcDecl); f != nil {
				failures = append(failures, *f)
			}
		case strings.HasPrefix(name, "Benchmark"):
			if f := checkBenchmarkFunc(funcDecl); f != nil {
				failures = append(failures, *f)
			}
		case strings.HasPrefix(name, "Fuzz"):
			if f := checkFuzzFunc(funcDecl); f != nil {
				failures = append(failures, *f)
			}
		case strings.HasPrefix(name, "Example"):
			if f := checkExampleFunc(funcDecl, topLevelIdents); f != nil {
				failures = append(failures, *f)
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoMalformedTestFunctionRule) Name() string {
	return "noMalformedTestFunction"
}

// Group returns the rule group.
func (*NoMalformedTestFunctionRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMalformedTestFunctionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// checkTestFunc validates a Test function has the correct signature.
// Test functions must have exactly one parameter of type *testing.T.
// The name after "Test" must either be empty (for TestMain) or start with an
// uppercase letter or underscore.
func checkTestFunc(fn *ast.FuncDecl) *lint.Failure {
	name := fn.Name.Name

	// "Test" alone is valid; names like "Testfoo" are not.
	if suffix := strings.TrimPrefix(name, "Test"); suffix != "" {
		if !isExportedOrUnderscore(suffix) {
			return &lint.Failure{
				Category:   lint.FailureCategoryNaming,
				Confidence: 1,
				Node:       fn,
				Failure:    fmt.Sprintf("test function %s has malformed name: first letter after 'Test' must be uppercase", name),
			}
		}
	}

	if !hasExactParam(fn, "testing", "T") {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       fn,
			Failure:    fmt.Sprintf("test function %s has wrong signature: must be func %s(t *testing.T)", name, name),
		}
	}

	return nil
}

// checkBenchmarkFunc validates a Benchmark function has the correct signature.
func checkBenchmarkFunc(fn *ast.FuncDecl) *lint.Failure {
	name := fn.Name.Name

	if suffix := strings.TrimPrefix(name, "Benchmark"); suffix != "" {
		if !isExportedOrUnderscore(suffix) {
			return &lint.Failure{
				Category:   lint.FailureCategoryNaming,
				Confidence: 1,
				Node:       fn,
				Failure:    fmt.Sprintf("benchmark function %s has malformed name: first letter after 'Benchmark' must be uppercase", name),
			}
		}
	}

	if !hasExactParam(fn, "testing", "B") {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       fn,
			Failure:    fmt.Sprintf("benchmark function %s has wrong signature: must be func %s(b *testing.B)", name, name),
		}
	}

	return nil
}

// checkFuzzFunc validates a Fuzz function has the correct signature.
func checkFuzzFunc(fn *ast.FuncDecl) *lint.Failure {
	name := fn.Name.Name

	if suffix := strings.TrimPrefix(name, "Fuzz"); suffix != "" {
		if !isExportedOrUnderscore(suffix) {
			return &lint.Failure{
				Category:   lint.FailureCategoryNaming,
				Confidence: 1,
				Node:       fn,
				Failure:    fmt.Sprintf("fuzz function %s has malformed name: first letter after 'Fuzz' must be uppercase", name),
			}
		}
	}

	if !hasExactParam(fn, "testing", "F") {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       fn,
			Failure:    fmt.Sprintf("fuzz function %s has wrong signature: must be func %s(f *testing.F)", name, name),
		}
	}

	return nil
}

// checkExampleFunc validates an Example function. Example functions must have
// no parameters and no results. If the name references a specific identifier
// (e.g. ExampleFoo), that identifier should exist as a top-level declaration.
func checkExampleFunc(fn *ast.FuncDecl, topLevelIdents map[string]bool) *lint.Failure {
	name := fn.Name.Name

	// Example functions must take no parameters and return no results.
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       fn,
			Failure:    fmt.Sprintf("example function %s must have no parameters", name),
		}
	}

	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       fn,
			Failure:    fmt.Sprintf("example function %s must have no return values", name),
		}
	}

	// Check if the suffix references an existing identifier.
	suffix := strings.TrimPrefix(name, "Example")
	if suffix == "" {
		// Just "Example" is valid — it's a whole-package example.
		return nil
	}

	// The suffix might contain an underscore for method examples: ExampleType_Method
	// Extract just the identifier part (before the first underscore).
	ident := suffix
	if idx := strings.Index(suffix, "_"); idx >= 0 {
		ident = suffix[:idx]
	}

	if ident == "" {
		// Example_ is valid (package example with suffix).
		return nil
	}

	// Check if the identifier exists in the file's top-level declarations.
	if !topLevelIdents[ident] {
		return &lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 0.8,
			Node:       fn,
			Failure:    fmt.Sprintf("example function %s references non-existent identifier %s", name, ident),
		}
	}

	return nil
}

// hasExactParam checks that a function has exactly one parameter of type
// *<pkg>.<typeName> (e.g. *testing.T).
func hasExactParam(fn *ast.FuncDecl, pkg, typeName string) bool {
	params := fn.Type.Params
	if params == nil || len(params.List) != 1 {
		return false
	}

	param := params.List[0]

	// The parameter must be a pointer: *testing.T
	starExpr, ok := param.Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	// The pointed-to type must be a selector: testing.T
	selExpr, ok := starExpr.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent, ok := selExpr.X.(*ast.Ident)
	if !ok {
		return false
	}

	return pkgIdent.Name == pkg && selExpr.Sel.Name == typeName
}

// isExportedOrUnderscore checks if the first character of s is an uppercase
// letter or an underscore.
func isExportedOrUnderscore(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return r == '_' || unicode.IsUpper(r)
}

// collectTopLevelIdents collects all top-level identifiers (functions, types,
// variables, constants) declared in the file.
func collectTopLevelIdents(f *ast.File) map[string]bool {
	idents := make(map[string]bool)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Recv != nil {
				// For methods, add the receiver type as an identifier.
				addReceiverType(idents, d.Recv)
			}
			idents[d.Name.Name] = true
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					idents[s.Name.Name] = true
				case *ast.ValueSpec:
					for _, n := range s.Names {
						idents[n.Name] = true
					}
				}
			}
		}
	}
	return idents
}

// addReceiverType extracts the type name from a method receiver and adds it
// to the identifiers map.
func addReceiverType(idents map[string]bool, recv *ast.FieldList) {
	if recv == nil || len(recv.List) == 0 {
		return
	}
	typ := recv.List[0].Type
	// Handle pointer receiver: *T
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	if ident, ok := typ.(*ast.Ident); ok {
		idents[ident.Name] = true
	}
}
