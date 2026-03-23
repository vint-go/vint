package no_zip_slip

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoZipSlipRule detects file path traversal when extracting zip or tar archives (Zip Slip vulnerability).
type NoZipSlipRule struct{}

// Apply applies the rule to given file.
func (r *NoZipSlipRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !hasArchiveImport(file.AST) {
		return nil
	}

	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkFuncBody(funcDecl.Body, func(f lint.Failure) {
			failures = append(failures, f)
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoZipSlipRule) Name() string {
	return "noZipSlip"
}

// Group returns the rule group.
func (*NoZipSlipRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoZipSlipRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// hasArchiveImport checks if the file imports archive/zip or archive/tar.
func hasArchiveImport(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}
		path := imp.Path.Value
		if path == `"archive/zip"` || path == `"archive/tar"` {
			return true
		}
	}
	return false
}

// checkFuncBody examines a function body for Zip Slip patterns.
// It looks for filepath.Join calls whose results are used in file operations
// without strings.HasPrefix validation.
func checkFuncBody(body *ast.BlockStmt, onFailure func(lint.Failure)) {
	// Collect all filepath.Join assignment targets
	type joinInfo struct {
		varName  string
		joinNode ast.Node
	}
	var joinVars []joinInfo

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, rhs := range assign.Rhs {
			if isFilepathJoinCall(rhs) {
				if i < len(assign.Lhs) {
					if ident, ok := assign.Lhs[i].(*ast.Ident); ok && ident.Name != "_" {
						joinVars = append(joinVars, joinInfo{
							varName:  ident.Name,
							joinNode: assign,
						})
					}
				}
			}
		}
		return true
	})

	// For each filepath.Join result, check if it is validated with strings.HasPrefix
	for _, jv := range joinVars {
		if !hasHasPrefixCheck(body, jv.varName) {
			onFailure(lint.Failure{
				Confidence: 1,
				Node:       jv.joinNode,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "potential Zip Slip: archive entry path not validated after filepath.Join",
			})
		}
	}
}

// isFilepathJoinCall checks if an expression is a call to filepath.Join.
func isFilepathJoinCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "filepath" && sel.Sel.Name == "Join"
}

// hasHasPrefixCheck checks if the function body contains a strings.HasPrefix
// call that includes the given variable name as an argument.
func hasHasPrefixCheck(body ast.Node, varName string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Name != "strings" || sel.Sel.Name != "HasPrefix" {
			return true
		}
		// Check if any argument references the variable
		for _, arg := range call.Args {
			if argIdent, ok := arg.(*ast.Ident); ok && argIdent.Name == varName {
				found = true
				return false
			}
		}
		return true
	})
	return found
}
