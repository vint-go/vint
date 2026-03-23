package no_capitalized_local

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoCapitalizedLocalRule detects capitalized names for local variables,
// function parameters, and return values.
type NoCapitalizedLocalRule struct {
	paramsOnly bool
	configured bool
}

// Configure validates and applies the rule configuration.
func (r *NoCapitalizedLocalRule) Configure(arguments lint.Arguments) error {
	r.paramsOnly = true // default
	r.configured = true

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noCapitalizedLocal" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		if normalizeOption(k) == normalizeOption("paramsOnly") {
			val, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for paramsOnly in "noCapitalizedLocal" rule; need bool but got %T`, v)
			}
			r.paramsOnly = val
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoCapitalizedLocalRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.paramsOnly = true
		r.configured = true
	}

	var failures []lint.Failure

	w := &lintCapitalizedLocal{
		paramsOnly: r.paramsOnly,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoCapitalizedLocalRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	if !r.configured {
		r.paramsOnly = true
		r.configured = true
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintCapitalizedLocal{
		paramsOnly: r.paramsOnly,
		onFailure:  onFailure,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoCapitalizedLocalRule) Name() string {
	return "noCapitalizedLocal"
}

// Group returns the rule group.
func (*NoCapitalizedLocalRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoCapitalizedLocalRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCapitalizedLocal struct {
	paramsOnly bool
	onFailure  func(lint.Failure)
}

func (w *lintCapitalizedLocal) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkFuncParams(n.Type)
		if !w.paramsOnly && n.Body != nil {
			w.checkLocalVars(n.Body)
		}
	case *ast.FuncLit:
		w.checkFuncParams(n.Type)
		if !w.paramsOnly && n.Body != nil {
			w.checkLocalVars(n.Body)
		}
	}
	return w
}

func (w *lintCapitalizedLocal) checkFuncParams(funcType *ast.FuncType) {
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			for _, name := range field.Names {
				w.checkName(name)
			}
		}
	}
	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			for _, name := range field.Names {
				w.checkName(name)
			}
		}
	}
}

func (w *lintCapitalizedLocal) checkLocalVars(body *ast.BlockStmt) {
	ast.Inspect(body, func(n ast.Node) bool {
		// Skip nested function literals; they will be handled by the outer walker.
		if _, ok := n.(*ast.FuncLit); ok {
			return false
		}

		switch stmt := n.(type) {
		case *ast.AssignStmt:
			if stmt.Tok == token.DEFINE {
				for _, lhs := range stmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						w.checkName(ident)
					}
				}
			}
		case *ast.DeclStmt:
			if genDecl, ok := stmt.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
				for _, spec := range genDecl.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range vs.Names {
							w.checkName(name)
						}
					}
				}
			}
		case *ast.RangeStmt:
			if stmt.Tok == token.DEFINE {
				if key, ok := stmt.Key.(*ast.Ident); ok {
					w.checkName(key)
				}
				if val, ok := stmt.Value.(*ast.Ident); ok {
					w.checkName(val)
				}
			}
		}
		return true
	})
}

func (w *lintCapitalizedLocal) checkName(ident *ast.Ident) {
	name := ident.Name
	if name == "_" || name == "" {
		return
	}

	firstRune, _ := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(firstRune) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Category:   lint.FailureCategoryNaming,
			Node:       ident,
			Failure:    fmt.Sprintf("local variable %s should not be capitalized", name),
		})
	}
}

func normalizeOption(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}
