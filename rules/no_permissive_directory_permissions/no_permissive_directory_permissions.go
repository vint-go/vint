package no_permissive_directory_permissions

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoPermissiveDirectoryPermissionsRule detects overly permissive directory permissions
// in calls to os.Mkdir and os.MkdirAll.
type NoPermissiveDirectoryPermissionsRule struct{}

// Apply applies the rule to given file.
func (r *NoPermissiveDirectoryPermissionsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintPermissiveDirPerms{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPermissiveDirectoryPermissionsRule) Name() string {
	return "noPermissiveDirectoryPermissions"
}

// Group returns the rule group.
func (*NoPermissiveDirectoryPermissionsRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoPermissiveDirectoryPermissionsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

const maxPermission = 0o750

type lintPermissiveDirPerms struct {
	onFailure func(lint.Failure)
}

func (w *lintPermissiveDirPerms) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for os.Mkdir or os.MkdirAll
	if !astutils.IsPkgDotName(ce.Fun, "os", "Mkdir") &&
		!astutils.IsPkgDotName(ce.Fun, "os", "MkdirAll") {
		return w
	}

	// os.Mkdir(name, perm) - perm is the second argument (index 1)
	// os.MkdirAll(path, perm) - perm is the second argument (index 1)
	if len(ce.Args) < 2 {
		return w
	}

	permArg := ce.Args[1]
	lit, ok := permArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return w
	}

	perm, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return w
	}

	if perm > maxPermission {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("directory permission %#o is more permissive than %#o", perm, maxPermission),
		})
	}

	return w
}
