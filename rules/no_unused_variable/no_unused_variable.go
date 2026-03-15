package no_unused_variable

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedVariableRule detects package-level variables that are declared but never
// read or referenced in the codebase.
type NoUnusedVariableRule struct {
	postStatementsAreReads bool
	localVariablesAreUsed  bool
	generatedIsUsed        bool
}

const (
	defaultPostStatementsAreReads = true
	defaultLocalVariablesAreUsed  = true
	defaultGeneratedIsUsed        = true
)

// Configure validates and applies the rule configuration.
func (r *NoUnusedVariableRule) Configure(arguments lint.Arguments) error {
	r.postStatementsAreReads = defaultPostStatementsAreReads
	r.localVariablesAreUsed = defaultLocalVariablesAreUsed
	r.generatedIsUsed = defaultGeneratedIsUsed

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnusedVariable" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch strings.ReplaceAll(k, "_", "-") {
		case "post-statements-are-reads":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "post-statements-are-reads" in "noUnusedVariable"; need bool but got %T`, v)
			}
			r.postStatementsAreReads = b
		case "local-variables-are-used":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "local-variables-are-used" in "noUnusedVariable"; need bool but got %T`, v)
			}
			r.localVariablesAreUsed = b
		case "generated-is-used":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "generated-is-used" in "noUnusedVariable"; need bool but got %T`, v)
			}
			r.generatedIsUsed = b
		}
	}

	return nil
}

// pkgVarInfo holds information about a package-level variable declaration.
type pkgVarInfo struct {
	name     string
	node     ast.Node // the *ast.Ident of the variable name
	fileName string
}

// Apply applies the rule to the given file.
func (r *NoUnusedVariableRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	pkg := file.Pkg

	// Phase 1: Collect all package-level variable declarations across all files.
	allVars := map[string]*pkgVarInfo{}
	entryVars := map[string]bool{} // variables considered always used

	for fname, f := range pkg.Files() {
		// If generated-is-used is enabled and file is generated, skip its vars.
		if r.generatedIsUsed && isGeneratedFile(f) {
			for _, decl := range f.AST.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.VAR {
					continue
				}
				for _, spec := range genDecl.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range vs.Names {
						entryVars[name.Name] = true
					}
				}
			}
			continue
		}

		for _, decl := range f.AST.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, nameIdent := range vs.Names {
					name := nameIdent.Name

					// Blank identifier is always considered used.
					if name == "_" {
						continue
					}

					// Exported variables are always considered used.
					if ast.IsExported(name) {
						entryVars[name] = true
						continue
					}

					allVars[name] = &pkgVarInfo{
						name:     name,
						node:     nameIdent,
						fileName: fname,
					}
				}
			}
		}
	}

	if len(allVars) == 0 {
		return nil
	}

	// Phase 2: Scan all files for reads/references to the collected variables.
	readVars := map[string]bool{}

	for _, f := range pkg.Files() {
		isTest := f.IsTest()
		for _, decl := range f.AST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Body == nil {
					continue
				}
				r.scanForReads(d.Body, allVars, readVars, isTest)
			case *ast.GenDecl:
				// Scan variable initializers and type expressions.
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						// Scan the values (initializers) for references to other vars.
						for _, val := range s.Values {
							r.scanExprForReads(val, allVars, readVars, s.Names)
						}
						// Scan the type expression if present.
						if s.Type != nil {
							r.scanExprForReads(s.Type, allVars, readVars, nil)
						}
					case *ast.TypeSpec:
						r.scanExprForReads(s.Type, allVars, readVars, nil)
					}
				}
			}
		}
	}

	// Phase 3: Report unused variables declared in this file.
	thisFileName := fileNameFromFile(file)

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, nameIdent := range vs.Names {
				name := nameIdent.Name
				if name == "_" {
					continue
				}
				if ast.IsExported(name) {
					continue
				}
				if entryVars[name] {
					continue
				}

				info, exists := allVars[name]
				if !exists {
					continue
				}
				// Only report if declared in this file.
				if info.fileName != thisFileName {
					continue
				}

				if readVars[name] {
					continue
				}

				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       nameIdent,
					Failure:    fmt.Sprintf("var %s is unused", name),
				})
			}
		}
	}

	return failures
}

// scanForReads walks a statement/block looking for identifier reads that reference
// package-level variables.
func (r *NoUnusedVariableRule) scanForReads(node ast.Node, allVars map[string]*pkgVarInfo, readVars map[string]bool, isTestFile bool) {
	// First, collect all LHS identifiers from simple assignments and inc/dec stmts
	// so we can distinguish writes from reads.
	writeOnlyIdents := map[*ast.Ident]bool{}
	postStmtIdents := map[*ast.Ident]bool{}

	ast.Inspect(node, func(n ast.Node) bool {
		switch s := n.(type) {
		case *ast.AssignStmt:
			if s.Tok == token.ASSIGN || s.Tok == token.DEFINE {
				for _, lhs := range s.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						if _, isPkgVar := allVars[ident.Name]; isPkgVar {
							writeOnlyIdents[ident] = true
						}
					}
				}
			} else {
				// Compound assignments (+=, -=, etc.) - these are both read and write.
				for _, lhs := range s.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						if _, isPkgVar := allVars[ident.Name]; isPkgVar {
							postStmtIdents[ident] = true
						}
					}
				}
			}
		case *ast.IncDecStmt:
			if ident, ok := s.X.(*ast.Ident); ok {
				if _, isPkgVar := allVars[ident.Name]; isPkgVar {
					postStmtIdents[ident] = true
				}
			}
		}
		return true
	})

	// Now walk again and mark all identifiers that reference package-level vars
	// as reads, except those that are pure writes (unless options say otherwise).
	ast.Inspect(node, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if _, isPkgVar := allVars[ident.Name]; !isPkgVar {
			return true
		}

		// Check if this is a write-only ident.
		if writeOnlyIdents[ident] {
			if isTestFile || r.postStatementsAreReads {
				readVars[ident.Name] = true
			}
			return true
		}

		// Check if this is a post-statement (++/--, +=, etc.) ident.
		if postStmtIdents[ident] {
			if isTestFile || r.postStatementsAreReads {
				readVars[ident.Name] = true
			}
			return true
		}

		// This is a read reference.
		readVars[ident.Name] = true
		return true
	})
}

// scanExprForReads scans an expression for reads of package-level variables.
// skipNames is a set of ident names to skip (e.g. when scanning a var initializer,
// skip the var being declared itself).
func (r *NoUnusedVariableRule) scanExprForReads(expr ast.Expr, allVars map[string]*pkgVarInfo, readVars map[string]bool, skipNames []*ast.Ident) {
	skipSet := map[string]bool{}
	for _, n := range skipNames {
		skipSet[n.Name] = true
	}

	ast.Inspect(expr, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		name := ident.Name
		if skipSet[name] {
			return true
		}
		if _, isPkgVar := allVars[name]; isPkgVar {
			readVars[name] = true
		}
		return true
	})
}

// fileNameFromFile extracts the file name that matches pkg.Files() keys.
func fileNameFromFile(file *lint.File) string {
	for name, f := range file.Pkg.Files() {
		if f == file {
			return name
		}
	}
	return ""
}

// isGeneratedFile checks if a file has a "Code generated" comment indicating
// it was automatically generated.
func isGeneratedFile(file *lint.File) bool {
	for _, cg := range file.AST.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, "Code generated") && strings.Contains(c.Text, "DO NOT EDIT") {
				return true
			}
		}
	}
	return false
}

// Name returns the rule name.
func (*NoUnusedVariableRule) Name() string {
	return "noUnusedVariable"
}

// Group returns the rule group.
func (*NoUnusedVariableRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedVariableRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
