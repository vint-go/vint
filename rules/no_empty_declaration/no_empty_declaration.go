package no_empty_declaration

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoEmptyDeclarationRule detects empty var, const, type, or import declaration blocks.
type NoEmptyDeclarationRule struct{}

// Apply applies the rule to given file.
func (r *NoEmptyDeclarationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintEmptyDecl{onFailure: onFailure}

	for _, decl := range file.AST.Decls {
		w.Visit(decl)
	}

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoEmptyDeclarationRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintEmptyDecl{onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*NoEmptyDeclarationRule) Name() string {
	return "noEmptyDeclaration"
}

// Group returns the rule group.
func (*NoEmptyDeclarationRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoEmptyDeclarationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintEmptyDecl struct {
	onFailure func(lint.Failure)
}

func (w *lintEmptyDecl) Visit(node ast.Node) ast.Visitor {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return w
	}

	// Only check parenthesized declaration blocks
	if genDecl.Lparen == token.NoPos {
		return w
	}

	// Check if the block is empty (no specs)
	if len(genDecl.Specs) > 0 {
		return w
	}

	var keyword string
	switch genDecl.Tok {
	case token.VAR:
		keyword = "var"
	case token.CONST:
		keyword = "const"
	case token.TYPE:
		keyword = "type"
	case token.IMPORT:
		keyword = "import"
	default:
		return w
	}

	w.onFailure(lint.Failure{
		Category:   lint.FailureCategoryStyle,
		Confidence: 1,
		Node:       genDecl,
		Failure:    fmt.Sprintf("empty %s declaration block should be removed", keyword),
	})

	return w
}
