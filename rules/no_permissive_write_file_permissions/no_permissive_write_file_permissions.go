package no_permissive_write_file_permissions

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoPermissiveWriteFilePermissionsRule detects overly permissive file permissions
// in calls to os.WriteFile and ioutil.WriteFile.
type NoPermissiveWriteFilePermissionsRule struct {
	maxPermission int64
}

// Configure implements lint.ConfigurableRule.
func (r *NoPermissiveWriteFilePermissionsRule) Configure(arguments lint.Arguments) error {
	r.maxPermission = defaultMaxWriteFilePermission

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noPermissiveWriteFilePermissions" rule, expecting a k,v map, got %T`, arguments[0])
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
					return fmt.Errorf(`invalid maxPermission in "noPermissiveWriteFilePermissions" rule: %w`, err)
				}
				r.maxPermission = parsed
			default:
				return fmt.Errorf(`invalid configuration value for maxPermission in "noPermissiveWriteFilePermissions" rule; need int64 or string but got %T`, v)
			}
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoPermissiveWriteFilePermissionsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	maxPerm := r.maxPermission
	if maxPerm == 0 {
		maxPerm = defaultMaxWriteFilePermission
	}
	w := &lintPermissiveWriteFilePerms{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		maxPermission: maxPerm,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPermissiveWriteFilePermissionsRule) Name() string {
	return "noPermissiveWriteFilePermissions"
}

// Group returns the rule group.
func (*NoPermissiveWriteFilePermissionsRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoPermissiveWriteFilePermissionsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

const defaultMaxWriteFilePermission = 0o600

type lintPermissiveWriteFilePerms struct {
	onFailure     func(lint.Failure)
	maxPermission int64
}

func (w *lintPermissiveWriteFilePerms) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check for os.WriteFile or ioutil.WriteFile
	isOsWriteFile := astutils.IsPkgDotName(ce.Fun, "os", "WriteFile")
	isIoutilWriteFile := astutils.IsPkgDotName(ce.Fun, "ioutil", "WriteFile")

	if !isOsWriteFile && !isIoutilWriteFile {
		return w
	}

	// os.WriteFile(name, data, perm) - perm is the third argument (index 2)
	// ioutil.WriteFile(name, data, perm) - perm is the third argument (index 2)
	if len(ce.Args) < 3 {
		return w
	}

	permArg := ce.Args[2]
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
