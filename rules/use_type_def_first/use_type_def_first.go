package use_type_def_first

import (
	"fmt"
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseTypeDefFirstRule detects method declarations preceding the type definition itself.
type UseTypeDefFirstRule struct{}

// Apply applies the rule to given file.
func (r *UseTypeDefFirstRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect all type names and the line where they are defined.
	typeDefLines := map[string]int{}

	// First pass: record all type definitions and their positions.
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			typeDefLines[typeSpec.Name.Name] = file.ToPosition(typeSpec.Pos()).Line
		}
	}

	// Second pass: check function declarations with receivers.
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if funcDecl.Recv == nil || len(funcDecl.Recv.List) == 0 {
			continue
		}

		// Extract the receiver type name.
		recvTypeName := receiverTypeName(funcDecl.Recv.List[0].Type)
		if recvTypeName == "" {
			continue
		}

		funcLine := file.ToPosition(funcDecl.Pos()).Line
		typeDefLine, exists := typeDefLines[recvTypeName]
		if !exists {
			// Type not defined in this file; skip.
			continue
		}

		if funcLine < typeDefLine {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryStyle,
				Node:       funcDecl,
				Failure:    fmt.Sprintf("method %s declared before type definition of %s", funcDecl.Name.Name, recvTypeName),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*UseTypeDefFirstRule) Name() string {
	return "useTypeDefFirst"
}

// Group returns the rule group.
func (*UseTypeDefFirstRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseTypeDefFirstRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// receiverTypeName extracts the type name from a receiver expression.
// It handles both value receivers (T) and pointer receivers (*T).
func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.IndexExpr:
		// Generic type: T[P]
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	case *ast.IndexListExpr:
		// Generic type: T[P1, P2]
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}
