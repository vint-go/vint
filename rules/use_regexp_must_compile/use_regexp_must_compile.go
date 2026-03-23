package use_regexp_must_compile

import (
	"go/ast"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseRegexpMustCompileRule detects regexp.Compile* calls that can be replaced
// with regexp.MustCompile*. When a regular expression is compiled with a constant
// pattern at package level or in an init function, MustCompile is preferred
// because it panics on invalid patterns, catching errors at startup rather than
// at runtime.
type UseRegexpMustCompileRule struct{}

// Apply applies the rule to the given file.
func (r *UseRegexpMustCompileRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	// Check package-level variable declarations
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				checkCompileCall(val, onFailure)
			}
		}
	}

	// Check init functions
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Name.Name != "init" || funcDecl.Recv != nil {
			continue
		}
		if funcDecl.Body == nil {
			continue
		}
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			checkCompileCall(n, onFailure)
			return true
		})
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseRegexpMustCompileRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, val := range vs.Values {
				checkCompileCall(val, onFailure)
			}
		}
	case *ast.FuncDecl:
		if n.Name.Name == "init" && n.Recv == nil && n.Body != nil {
			ast.Inspect(n.Body, func(nn ast.Node) bool {
				if nn == nil {
					return false
				}
				checkCompileCall(nn, onFailure)
				return true
			})
		}
	}

	return failures
}

// checkCompileCall checks if a node is a regexp.Compile or regexp.CompilePOSIX
// call with a constant string pattern, and reports a failure if so.
func checkCompileCall(node ast.Node, onFailure func(lint.Failure)) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return
	}

	var compileFunc string
	if astutils.IsPkgDotName(call.Fun, "regexp", "Compile") {
		compileFunc = "regexp.Compile"
	} else if astutils.IsPkgDotName(call.Fun, "regexp", "CompilePOSIX") {
		compileFunc = "regexp.CompilePOSIX"
	} else {
		return
	}

	// Must have at least one argument (the pattern)
	if len(call.Args) == 0 {
		return
	}

	// Only flag when the pattern argument is a string literal
	if !astutils.IsStringLiteral(call.Args[0]) {
		return
	}

	mustFunc := strings.Replace(compileFunc, "Compile", "MustCompile", 1)

	onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryStyle,
		Failure:    compileFunc + " can be replaced by " + mustFunc,
	})
}

// Name returns the rule name.
func (*UseRegexpMustCompileRule) Name() string {
	return "useRegexpMustCompile"
}

// Group returns the rule group.
func (*UseRegexpMustCompileRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseRegexpMustCompileRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
