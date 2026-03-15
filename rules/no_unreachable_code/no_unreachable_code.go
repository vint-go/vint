package no_unreachable_code

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnreachableCodeRule detects code that can never be executed because it appears
// after a statement that unconditionally exits the function, such as return, panic,
// os.Exit, log.Fatal, or an infinite loop with no break.
type NoUnreachableCodeRule struct{}

// Apply applies the rule to given file.
func (r *NoUnreachableCodeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	testingFunctions := map[string]bool{
		"Fatal":   true,
		"Fatalf":  true,
		"FailNow": true,
	}
	branchingFunctions := map[string]map[string]bool{
		"os": {"Exit": true},
		"log": {
			"Fatal":   true,
			"Fatalf":  true,
			"Fatalln": true,
			"Panic":   true,
			"Panicf":  true,
			"Panicln": true,
		},
		"t": testingFunctions,
		"b": testingFunctions,
		"f": testingFunctions,
	}

	w := &lintUnreachableCode{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
		branchingFunctions: branchingFunctions,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoUnreachableCodeRule) Name() string {
	return "noUnreachableCode"
}

// Group returns the rule group.
func (*NoUnreachableCodeRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnreachableCodeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnreachableCode struct {
	onFailure          func(lint.Failure)
	branchingFunctions map[string]map[string]bool
}

func (w *lintUnreachableCode) Visit(node ast.Node) ast.Visitor {
	blk, ok := node.(*ast.BlockStmt)
	if !ok {
		return w
	}

	if len(blk.List) < 2 {
		return w
	}

loop:
	for i, stmt := range blk.List[:len(blk.List)-1] {
		next := blk.List[i+1]
		if _, ok := next.(*ast.LabeledStmt); ok {
			continue // skip if next statement is labeled
		}

		switch s := stmt.(type) {
		case *ast.ReturnStmt:
			w.onFailure(newUnreachableCodeFailure(s))
			break loop
		case *ast.BranchStmt:
			token := s.Tok.String()
			if token != "fallthrough" {
				w.onFailure(newUnreachableCodeFailure(s))
				break loop
			}
		case *ast.ExprStmt:
			ce, ok := s.X.(*ast.CallExpr)
			if !ok {
				continue
			}

			// Check for panic() built-in
			if ident, ok := ce.Fun.(*ast.Ident); ok && ident.Name == "panic" {
				if _, ok := next.(*ast.ReturnStmt); ok {
					continue // return after panic is allowed to satisfy function signature
				}
				w.onFailure(newUnreachableCodeFailure(s))
				break loop
			}

			// Check for pkg.Func() calls like os.Exit, log.Fatal, etc.
			fc, ok := ce.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}

			id, ok := fc.X.(*ast.Ident)
			if !ok {
				continue
			}
			fn := fc.Sel.Name
			pkg := id.Name
			if !w.branchingFunctions[pkg][fn] {
				continue
			}

			if _, ok := next.(*ast.ReturnStmt); ok {
				continue // return after branching function is allowed to satisfy function signature
			}

			w.onFailure(newUnreachableCodeFailure(s))
			break loop
		case *ast.ForStmt:
			// Detect infinite loops: for {} or for { ... } with no break/goto/return inside
			if s.Cond == nil && !containsBreakOrGoto(s.Body) {
				if _, ok := next.(*ast.ReturnStmt); ok {
					continue // return after infinite loop is allowed to satisfy function signature
				}
				w.onFailure(newUnreachableCodeFailure(s))
				break loop
			}
		}
	}

	return w
}

// containsBreakOrGoto checks if a block statement contains a break, goto,
// or return statement that would allow exiting the loop.
func containsBreakOrGoto(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}
	found := false
	ast.Inspect(block, func(n ast.Node) bool {
		if found {
			return false
		}
		switch s := n.(type) {
		case *ast.BranchStmt:
			if s.Tok.String() == "break" || s.Tok.String() == "goto" {
				found = true
				return false
			}
		case *ast.ReturnStmt:
			found = true
			return false
		case *ast.ForStmt, *ast.RangeStmt:
			// Don't look inside nested loops; a break there exits the inner loop
			if n != block {
				return false
			}
		case *ast.SelectStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			// Don't look inside select/switch; a break there exits the select/switch
			return false
		}
		return true
	})
	return found
}

func newUnreachableCodeFailure(node ast.Node) lint.Failure {
	return lint.Failure{
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryLogic,
		Failure:    "unreachable code after this statement",
	}
}
