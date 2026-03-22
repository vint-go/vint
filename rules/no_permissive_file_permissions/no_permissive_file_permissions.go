package no_permissive_file_permissions

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoPermissiveFilePermissionsRule detects overly permissive file permissions
// in calls to os.OpenFile and os.Chmod.
type NoPermissiveFilePermissionsRule struct {
	maxPermission int64
}

// Configure implements lint.ConfigurableRule.
func (r *NoPermissiveFilePermissionsRule) Configure(arguments lint.Arguments) error {
	r.maxPermission = defaultMaxFilePermission

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noPermissiveFilePermissions" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch k {
		case "maxPermission":
			switch val := v.(type) {
			case int64:
				r.maxPermission = val
			case string:
				parsed, err := strconv.ParseInt(val, 0, 64)
				if err != nil {
					return fmt.Errorf(`invalid maxPermission in "noPermissiveFilePermissions" rule: %w`, err)
				}
				r.maxPermission = parsed
			default:
				return fmt.Errorf(`invalid configuration value for maxPermission in "noPermissiveFilePermissions" rule; need int64 or string but got %T`, v)
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoPermissiveFilePermissionsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	maxPerm := r.maxPermission
	if maxPerm == 0 {
		maxPerm = defaultMaxFilePermission
	}
	w := &lintPermissiveFilePerms{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		maxPermission: maxPerm,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPermissiveFilePermissionsRule) Name() string {
	return "noPermissiveFilePermissions"
}

// Group returns the rule group.
func (*NoPermissiveFilePermissionsRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoPermissiveFilePermissionsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

const defaultMaxFilePermission = 0o600

type lintPermissiveFilePerms struct {
	onFailure     func(lint.Failure)
	maxPermission int64
}

func (w *lintPermissiveFilePerms) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for os.OpenFile or os.Chmod
	isOpenFile := astutils.IsPkgDotName(ce.Fun, "os", "OpenFile")
	isChmod := astutils.IsPkgDotName(ce.Fun, "os", "Chmod")

	if !isOpenFile && !isChmod {
		return w
	}

	// os.OpenFile(name, flag, perm) - perm is the third argument (index 2)
	// os.Chmod(name, mode) - mode is the second argument (index 1)
	var permArgIdx int
	if isOpenFile {
		permArgIdx = 2
		if len(ce.Args) < 3 {
			return w
		}
	} else {
		permArgIdx = 1
		if len(ce.Args) < 2 {
			return w
		}
	}

	permArg := ce.Args[permArgIdx]
	lit, ok := permArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return w
	}

	perm, err := strconv.ParseInt(lit.Value, 0, 64)
	if err != nil {
		return w
	}

	if perm > w.maxPermission {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("file permission %#o is more permissive than %#o", perm, w.maxPermission),
		})
	}

	return w
}
