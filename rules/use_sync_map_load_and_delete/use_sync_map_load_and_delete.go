package use_sync_map_load_and_delete

import (
	"go/ast"
	"go/types"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseSyncMapLoadAndDeleteRule detects sync.Map Load followed by Delete
// that can be replaced with the atomic LoadAndDelete method.
type UseSyncMapLoadAndDeleteRule struct{}

// Apply applies the rule to given file.
func (r *UseSyncMapLoadAndDeleteRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkBlockStmt(file.Pkg, funcDecl.Body, onFailure)
	}

	return failures
}

func checkBlockStmt(pkg *lint.Package, block *ast.BlockStmt, onFailure func(lint.Failure)) {
	stmts := block.List

	for i := 0; i < len(stmts)-1; i++ {
		// Look for: v, ok := m.Load(key) followed by m.Delete(key)
		assignStmt, ok := stmts[i].(*ast.AssignStmt)
		if !ok {
			// Also recurse into nested blocks
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		// The assignment should be from a call expression: m.Load(key)
		if len(assignStmt.Rhs) != 1 {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		loadCall, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		loadSel, ok := loadCall.Fun.(*ast.SelectorExpr)
		if !ok || loadSel.Sel.Name != "Load" {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		if len(loadCall.Args) != 1 {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		// Check that the receiver is a sync.Map
		if !isSyncMapType(pkg, loadSel.X) {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		loadReceiver := astutils.GoFmt(loadSel.X)
		loadKey := astutils.GoFmt(loadCall.Args[0])

		// Check if the next statement is m.Delete(key)
		nextStmt := stmts[i+1]
		exprStmt, ok := nextStmt.(*ast.ExprStmt)
		if !ok {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		deleteCall, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		deleteSel, ok := deleteCall.Fun.(*ast.SelectorExpr)
		if !ok || deleteSel.Sel.Name != "Delete" {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		if len(deleteCall.Args) != 1 {
			inspectNestedBlocks(pkg, stmts[i], onFailure)
			continue
		}

		deleteReceiver := astutils.GoFmt(deleteSel.X)
		deleteKey := astutils.GoFmt(deleteCall.Args[0])

		if loadReceiver == deleteReceiver && loadKey == deleteKey {
			onFailure(lint.Failure{
				Category:   lint.FailureCategoryOptimization,
				Confidence: 1,
				Node:       assignStmt,
				Failure:    "use " + loadReceiver + ".LoadAndDelete(" + loadKey + ") instead of separate Load and Delete calls",
			})
			// Skip the delete statement since it's part of the pattern
			i++
			continue
		}

		inspectNestedBlocks(pkg, stmts[i], onFailure)
	}

	// Handle the last statement's nested blocks
	if len(stmts) > 0 {
		inspectNestedBlocks(pkg, stmts[len(stmts)-1], onFailure)
	}
}

func inspectNestedBlocks(pkg *lint.Package, stmt ast.Stmt, onFailure func(lint.Failure)) {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
		if s.Else != nil {
			if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
				checkBlockStmt(pkg, elseBlock, onFailure)
			} else if elseIf, ok := s.Else.(*ast.IfStmt); ok {
				inspectNestedBlocks(pkg, elseIf, onFailure)
			}
		}
	case *ast.ForStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
	case *ast.RangeStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
	case *ast.SwitchStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
	case *ast.TypeSwitchStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
	case *ast.SelectStmt:
		checkBlockStmt(pkg, s.Body, onFailure)
	case *ast.CaseClause:
		// Process statements in case clauses
		block := &ast.BlockStmt{List: s.Body}
		checkBlockStmt(pkg, block, onFailure)
	case *ast.CommClause:
		block := &ast.BlockStmt{List: s.Body}
		checkBlockStmt(pkg, block, onFailure)
	case *ast.BlockStmt:
		checkBlockStmt(pkg, s, onFailure)
	}
}

func isSyncMapType(pkg *lint.Package, expr ast.Expr) bool {
	t := pkg.TypeOf(expr)
	if t == nil {
		return false
	}

	// Dereference pointer if needed
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	return obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == "sync" && obj.Name() == "Map"
}

// Name returns the rule name.
func (*UseSyncMapLoadAndDeleteRule) Name() string {
	return "useSyncMapLoadAndDelete"
}

// Group returns the rule group.
func (*UseSyncMapLoadAndDeleteRule) Group() string {
	return "performance"
}

// RequiresTypecheck returns true because this rule needs type information.
func (*UseSyncMapLoadAndDeleteRule) RequiresTypecheck() bool {
	return true
}

// CacheTier returns the cache tier for this rule.
func (*UseSyncMapLoadAndDeleteRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
