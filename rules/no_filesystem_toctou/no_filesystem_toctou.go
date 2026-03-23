package no_filesystem_toctou

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoFilesystemToctouRule detects filesystem TOCTOU (Time-of-Check-Time-of-Use)
// race conditions in filepath.Walk and filepath.WalkDir callbacks.
type NoFilesystemToctouRule struct{}

// Apply applies the rule to given file.
func (r *NoFilesystemToctouRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintNoFilesystemToctou{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoFilesystemToctouRule) Name() string {
	return "noFilesystemToctou"
}

// Group returns the rule group.
func (*NoFilesystemToctouRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoFilesystemToctouRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintNoFilesystemToctou struct {
	onFailure func(lint.Failure)
}

// fileOps maps package names to function names that perform file operations
// susceptible to TOCTOU when used with the path argument of Walk/WalkDir callbacks.
var fileOps = map[string]map[string]bool{
	"os": {
		"Open":      true,
		"OpenFile":  true,
		"ReadFile":  true,
		"Create":    true,
		"Remove":    true,
		"RemoveAll": true,
		"Stat":      true,
		"Lstat":     true,
		"Chmod":     true,
		"Chown":     true,
		"Rename":    true,
		"Link":      true,
		"Symlink":   true,
		"Mkdir":     true,
		"MkdirAll":  true,
	},
	"ioutil": {
		"ReadFile": true,
	},
}

func (w *lintNoFilesystemToctou) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a call to filepath.Walk or filepath.WalkDir
	isWalk := astutils.IsPkgDotName(ce.Fun, "filepath", "Walk")
	isWalkDir := astutils.IsPkgDotName(ce.Fun, "filepath", "WalkDir")
	if !isWalk && !isWalkDir {
		return w
	}

	// filepath.Walk(root, fn) and filepath.WalkDir(root, fn) both take 2 args
	if len(ce.Args) < 2 {
		return w
	}

	// The second argument is the walk function (callback)
	funcLit, ok := ce.Args[1].(*ast.FuncLit)
	if !ok {
		return w // callback is not an inline function literal, skip
	}

	// Get the name of the "path" parameter (first parameter of the callback)
	pathParamName := extractPathParamName(funcLit)
	if pathParamName == "" {
		return w
	}

	// Search the callback body for file operations that use the path parameter
	w.checkCallbackBody(funcLit.Body, pathParamName, isWalk)

	return w // continue walking other nodes
}

// extractPathParamName extracts the name of the first parameter (the path parameter)
// from a Walk/WalkDir callback function literal.
func extractPathParamName(funcLit *ast.FuncLit) string {
	if funcLit.Type == nil || funcLit.Type.Params == nil {
		return ""
	}
	params := funcLit.Type.Params.List
	if len(params) == 0 {
		return ""
	}
	// The first parameter is the path string
	firstParam := params[0]
	if len(firstParam.Names) == 0 {
		return ""
	}
	return firstParam.Names[0].Name
}

// checkCallbackBody walks the callback body looking for file operations
// that use the path parameter directly.
func (w *lintNoFilesystemToctou) checkCallbackBody(body *ast.BlockStmt, pathParamName string, isWalk bool) {
	if body == nil {
		return
	}

	walkFunc := "filepath.Walk"
	if !isWalk {
		walkFunc = "filepath.WalkDir"
	}

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		for pkg, funcs := range fileOps {
			for funcName := range funcs {
				if astutils.IsPkgDotName(call.Fun, pkg, funcName) {
					if w.usesPathParam(call.Args, pathParamName) {
						msg := fmt.Sprintf("potential TOCTOU race: %s.%s uses path from %s callback", pkg, funcName, walkFunc)
						w.onFailure(lint.Failure{
							Confidence: 1,
							Node:       call,
							Category:   lint.FailureCategoryBadPractice,
							Failure:    msg,
						})
					}
					return true
				}
			}
		}

		return true
	})
}

// usesPathParam checks if any of the arguments to a function call
// directly reference the path parameter variable.
func (w *lintNoFilesystemToctou) usesPathParam(args []ast.Expr, pathParamName string) bool {
	for _, arg := range args {
		if containsIdent(arg, pathParamName) {
			return true
		}
	}
	return false
}

// containsIdent checks if an expression contains a reference to the given identifier.
// It handles simple identifiers and also expressions like filepath.Join(path, ...).
func containsIdent(expr ast.Expr, name string) bool {
	found := false
	ast.Inspect(expr, func(n ast.Node) bool {
		if found {
			return false
		}
		if ident, ok := n.(*ast.Ident); ok && ident.Name == name {
			found = true
			return false
		}
		return true
	})
	return found
}
