package no_copied_lock

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoCopiedLockRule checks for locks erroneously passed by value.
type NoCopiedLockRule struct{}

// Apply applies the rule to given file.
func (r *NoCopiedLockRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	w := &lintNoCopiedLock{
		pkg: file.Pkg,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoCopiedLockRule) Name() string {
	return "noCopiedLock"
}

// Group returns the rule group.
func (*NoCopiedLockRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoCopiedLockRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoCopiedLockRule) RequiresTypecheck() bool {
	return true
}

type lintNoCopiedLock struct {
	pkg       *lint.Package
	onFailure func(lint.Failure)
}

func (w *lintNoCopiedLock) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFuncParams(n)
	case *ast.AssignStmt:
		w.checkAssignment(n)
	case *ast.ReturnStmt:
		w.checkReturn(n)
	case *ast.RangeStmt:
		w.checkRange(n)
	}
	return w
}

// checkFuncParams checks function parameters for lock values passed by value.
func (w *lintNoCopiedLock) checkFuncParams(fn *ast.FuncDecl) {
	if fn.Type.Params == nil {
		return
	}
	for _, field := range fn.Type.Params.List {
		t := w.pkg.TypeOf(field.Type)
		if t == nil {
			continue
		}
		if path := lockPath(t); path != "" {
			for _, name := range field.Names {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       field,
					Failure:    fmt.Sprintf("%s passes lock by value: %s", name.Name, path),
				})
			}
			if len(field.Names) == 0 {
				w.onFailure(lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       field,
					Failure:    fmt.Sprintf("passes lock by value: %s", path),
				})
			}
		}
	}
}

// checkAssignment checks assignments for copying lock values.
func (w *lintNoCopiedLock) checkAssignment(assign *ast.AssignStmt) {
	for _, rhs := range assign.Rhs {
		// Skip composite literals - they create new values, not copies.
		if _, ok := rhs.(*ast.CompositeLit); ok {
			continue
		}
		// Skip unary & expressions - taking address, not copying.
		if unary, ok := rhs.(*ast.UnaryExpr); ok && unary.Op.String() == "&" {
			continue
		}
		t := w.pkg.TypeOf(rhs)
		if t == nil {
			continue
		}
		if path := lockPath(t); path != "" {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       assign,
				Failure:    fmt.Sprintf("assignment copies lock value: %s", path),
			})
			return // one failure per assignment statement
		}
	}
}

// checkReturn checks return statements for copying lock values.
func (w *lintNoCopiedLock) checkReturn(ret *ast.ReturnStmt) {
	for _, result := range ret.Results {
		t := w.pkg.TypeOf(result)
		if t == nil {
			continue
		}
		if path := lockPath(t); path != "" {
			w.onFailure(lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       ret,
				Failure:    fmt.Sprintf("return copies lock value: %s", path),
			})
			return // one failure per return statement
		}
	}
}

// checkRange checks range loop variables for copying lock values.
func (w *lintNoCopiedLock) checkRange(rangeStmt *ast.RangeStmt) {
	if rangeStmt.Value == nil {
		return
	}
	t := w.pkg.TypeOf(rangeStmt.Value)
	if t == nil {
		return
	}
	if path := lockPath(t); path != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       rangeStmt,
			Failure:    fmt.Sprintf("range var copies lock value: %s", path),
		})
	}
}

// lockPath checks whether the given type contains a sync lock type.
// It returns a description of the path to the lock, or "" if no lock is found.
func lockPath(t types.Type) string {
	// Unwrap pointer types: pointers don't copy the underlying value.
	if _, ok := t.(*types.Pointer); ok {
		return ""
	}

	// Check if this is a named type from the sync package that implements Locker
	// or is a known non-copyable type.
	if named, ok := t.(*types.Named); ok {
		obj := named.Obj()
		if obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "sync" {
			switch obj.Name() {
			case "Mutex", "RWMutex", "WaitGroup", "Cond", "Once", "Pool", "Map":
				return fmt.Sprintf("sync.%s", obj.Name())
			}
		}
		// Also check if this implements sync.Locker (for custom types wrapping locks)
		// but for now, check the underlying struct fields.
		if st, ok := named.Underlying().(*types.Struct); ok {
			for i := 0; i < st.NumFields(); i++ {
				field := st.Field(i)
				if path := lockPath(field.Type()); path != "" {
					return fmt.Sprintf("%s.%s (%s)", named.Obj().Name(), field.Name(), path)
				}
			}
		}
	}

	// Check struct types directly (unnamed structs).
	if st, ok := t.(*types.Struct); ok {
		for i := 0; i < st.NumFields(); i++ {
			field := st.Field(i)
			if path := lockPath(field.Type()); path != "" {
				return path
			}
		}
	}

	// Check array types.
	if arr, ok := t.(*types.Array); ok {
		return lockPath(arr.Elem())
	}

	return ""
}
