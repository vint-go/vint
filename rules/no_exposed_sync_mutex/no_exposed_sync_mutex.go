package no_exposed_sync_mutex

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoExposedSyncMutexRule detects exported sync.Mutex and sync.RWMutex fields
// that inadvertently promote Lock/Unlock methods to a struct's public API.
type NoExposedSyncMutexRule struct{}

// Apply applies the rule to given file.
func (r *NoExposedSyncMutexRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintExposedMutex{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoExposedSyncMutexRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}
	w := &lintExposedMutex{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoExposedSyncMutexRule) Name() string {
	return "noExposedSyncMutex"
}

// Group returns the rule group.
func (*NoExposedSyncMutexRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoExposedSyncMutexRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintExposedMutex struct {
	onFailure func(lint.Failure)
}

func (w *lintExposedMutex) Visit(node ast.Node) ast.Visitor {
	structType, ok := node.(*ast.StructType)
	if !ok {
		return w
	}

	if structType.Fields == nil {
		return w
	}

	for _, field := range structType.Fields.List {
		if !isExposedMutexField(field) {
			continue
		}

		typeName := astutils.GoFmt(field.Type)
		if len(field.Names) == 0 {
			// Embedded field (unnamed): sync.Mutex or sync.RWMutex
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Node:       field,
				Failure:    fmt.Sprintf("embedded %s exposes Lock and Unlock methods to the public API; use an unexported field instead", typeName),
			})
		} else {
			// Named exported field
			for _, name := range field.Names {
				if name.IsExported() {
					w.onFailure(lint.Failure{
						Category:   lint.FailureCategoryStyle,
						Confidence: 1,
						Node:       field,
						Failure:    fmt.Sprintf("exported field %s of type %s exposes Lock and Unlock methods to the public API; use an unexported field instead", name.Name, typeName),
					})
				}
			}
		}
	}

	return w
}

// isExposedMutexField checks if a field is a sync.Mutex or sync.RWMutex
// that exposes methods to the public API (either unnamed embedded or exported named).
func isExposedMutexField(field *ast.Field) bool {
	sel, ok := field.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name != "sync" {
		return false
	}

	if sel.Sel.Name != "Mutex" && sel.Sel.Name != "RWMutex" {
		return false
	}

	// Unnamed embedded field (no names) always exposes methods
	if len(field.Names) == 0 {
		return true
	}

	// Named field: only a problem if the name is exported
	for _, name := range field.Names {
		if name.IsExported() {
			return true
		}
	}

	return false
}
