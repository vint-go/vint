package no_ssh_auth_bypass

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoSshAuthBypassRule detects stateful misuse of ssh.PublicKeyCallback
// that can lead to authentication bypass by setting Extensions in
// ssh.Permissions within the callback.
type NoSshAuthBypassRule struct{}

// Apply applies the rule to given file.
func (r *NoSshAuthBypassRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSshAuthBypass{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSshAuthBypassRule) Name() string {
	return "noSshAuthBypass"
}

// Group returns the rule group.
func (*NoSshAuthBypassRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSshAuthBypassRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSshAuthBypass struct {
	onFailure func(lint.Failure)
}

func (w *lintSshAuthBypass) Visit(node ast.Node) ast.Visitor {
	// Look for composite literals: &ssh.ServerConfig{...}
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return w
	}

	if !isSSHServerConfig(compositeLit.Type) {
		return w
	}

	// Find the PublicKeyCallback field
	for _, elt := range compositeLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		ident, ok := kv.Key.(*ast.Ident)
		if !ok || ident.Name != "PublicKeyCallback" {
			continue
		}

		// Check if the value is a function literal
		funcLit, ok := kv.Value.(*ast.FuncLit)
		if !ok {
			continue
		}

		// Check if the function body returns ssh.Permissions with Extensions set
		if containsPermissionsWithExtensions(funcLit.Body) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       kv,
				Category:   lint.FailureCategoryBadPractice,
				Failure:    "PublicKeyCallback sets Extensions in ssh.Permissions, which may lead to authentication bypass",
			})
		}
	}

	return w
}

// isSSHServerConfig checks if the expression is ssh.ServerConfig.
func isSSHServerConfig(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "ssh" && sel.Sel.Name == "ServerConfig"
}

// containsPermissionsWithExtensions checks if a function body contains
// a return of &ssh.Permissions{Extensions: ...} with a non-empty Extensions map.
func containsPermissionsWithExtensions(body *ast.BlockStmt) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if found {
			return false
		}

		// Look for composite literals of ssh.Permissions or &ssh.Permissions
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		if !isSSHPermissions(compLit.Type) {
			return true
		}

		// Check if Extensions field is set with a non-empty map
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			ident, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}
			if ident.Name == "Extensions" {
				// Check if the value is a non-empty composite literal (map)
				mapLit, ok := kv.Value.(*ast.CompositeLit)
				if ok && len(mapLit.Elts) > 0 {
					found = true
					return false
				}
			}
		}

		return true
	})
	return found
}

// isSSHPermissions checks if the expression is ssh.Permissions.
func isSSHPermissions(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "ssh" && sel.Sel.Name == "Permissions"
}
