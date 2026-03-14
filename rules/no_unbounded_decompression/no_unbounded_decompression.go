package no_unbounded_decompression

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnboundedDecompressionRule detects unbounded use of io.Copy with decompression readers,
// which can lead to decompression bomb (zip bomb) attacks.
type NoUnboundedDecompressionRule struct{}

// decompressionConstructors maps package alias to function names that create decompression readers.
var decompressionConstructors = map[string]string{
	"gzip":  "NewReader",
	"zlib":  "NewReader",
	"flate": "NewReader",
	"bzip2": "NewReader",
}

// Apply applies the rule to given file.
func (r *NoUnboundedDecompressionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkFunctionBody(funcDecl.Body, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnboundedDecompressionRule) Name() string {
	return "noUnboundedDecompression"
}

// Group returns the rule group.
func (*NoUnboundedDecompressionRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnboundedDecompressionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// checkFunctionBody analyzes a function body for unbounded decompression patterns.
func checkFunctionBody(body *ast.BlockStmt, onFailure func(lint.Failure)) {
	// Phase 1: Collect variable names assigned from decompression reader constructors
	decompVars := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Check if any RHS is a decompression constructor call.
		// For multi-value returns (e.g., gz, err := gzip.NewReader(src)),
		// only the first LHS variable is the reader.
		for i, rhs := range assign.Rhs {
			if !isDecompressionConstructorCall(rhs) {
				continue
			}

			if len(assign.Rhs) == 1 && len(assign.Lhs) > 1 {
				// Multi-value return: only the first LHS is the reader
				if ident, ok := assign.Lhs[0].(*ast.Ident); ok && ident.Name != "_" {
					decompVars[ident.Name] = true
				}
			} else if i < len(assign.Lhs) {
				// 1:1 assignment
				if ident, ok := assign.Lhs[i].(*ast.Ident); ok && ident.Name != "_" {
					decompVars[ident.Name] = true
				}
			}
		}
		return true
	})

	if len(decompVars) == 0 {
		return
	}

	// Phase 2: Find io.Copy calls where the source is a decompression reader variable
	ast.Inspect(body, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !astutils.IsPkgDotName(ce.Fun, "io", "Copy") {
			return true
		}

		if len(ce.Args) < 2 {
			return true
		}

		src := ce.Args[1]
		ident, ok := src.(*ast.Ident)
		if !ok {
			return true
		}

		if decompVars[ident.Name] {
			onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "use io.CopyN or io.LimitReader to prevent decompression bomb attacks",
			})
		}

		return true
	})
}

// isDecompressionConstructorCall checks if an expression is a call to a known
// decompression reader constructor (e.g., gzip.NewReader, zlib.NewReader).
func isDecompressionConstructorCall(expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	for pkg, funcName := range decompressionConstructors {
		if astutils.IsPkgDotName(ce.Fun, pkg, funcName) {
			return true
		}
	}

	return false
}
