package no_unused_constant

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedConstantRule detects constants that are declared but never
// referenced anywhere in the codebase.
type NoUnusedConstantRule struct {
	generatedIsUsed bool
}

const defaultGeneratedIsUsed = true

// Configure validates and applies the rule configuration.
func (r *NoUnusedConstantRule) Configure(arguments lint.Arguments) error {
	r.generatedIsUsed = defaultGeneratedIsUsed

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnusedConstant" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch strings.ReplaceAll(k, "_", "-") {
		case "generated-is-used":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "generated-is-used" in "noUnusedConstant"; need bool but got %T`, v)
			}
			r.generatedIsUsed = b
		}
	}

	return nil
}

// pkgConstInfo holds information about a package-level constant declaration.
type pkgConstInfo struct {
	name      string
	node      ast.Node // the *ast.Ident of the constant name
	fileName  string
	blockKeys []string // names of all constants in the same const block (for iota groups)
}

// Apply applies the rule to the given file.
func (r *NoUnusedConstantRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	pkg := file.Pkg

	// Phase 1: Collect all package-level constant declarations across all files.
	allConsts := map[string]*pkgConstInfo{}
	entryConsts := map[string]bool{} // constants considered always used

	// constBlockMap maps each constant name to its block group key list.
	// If one constant in a block is used, all are considered used (iota groups).
	constBlocks := map[string][]string{} // blockKey -> list of const names in that block

	for fname, f := range pkg.Files() {
		// If generated-is-used is enabled and file is generated, skip its consts.
		if r.generatedIsUsed && isGeneratedFile(f) {
			for _, decl := range f.AST.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.CONST {
					continue
				}
				for _, spec := range genDecl.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for _, name := range vs.Names {
						entryConsts[name.Name] = true
					}
				}
			}
			continue
		}

		for _, decl := range f.AST.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			// Collect all names in this const block for iota group handling.
			var blockNames []string
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				for _, nameIdent := range vs.Names {
					blockNames = append(blockNames, nameIdent.Name)
				}
			}

			// Determine a block key (use first name as block identifier).
			isBlock := len(genDecl.Specs) > 1
			if isBlock && len(blockNames) > 0 {
				constBlocks[blockNames[0]] = blockNames
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

					// Exported constants are always considered used.
					if ast.IsExported(name) {
						entryConsts[name] = true
						continue
					}

					var bk []string
					if isBlock {
						bk = blockNames
					}

					allConsts[name] = &pkgConstInfo{
						name:      name,
						node:      nameIdent,
						fileName:  fname,
						blockKeys: bk,
					}
				}
			}
		}
	}

	if len(allConsts) == 0 {
		return nil
	}

	// Phase 2: Scan all files for reads/references to the collected constants.
	readConsts := map[string]bool{}

	for _, f := range pkg.Files() {
		for _, decl := range f.AST.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				if d.Body == nil {
					continue
				}
				scanForConstReads(d.Body, allConsts, readConsts)
			case *ast.GenDecl:
				// Scan variable/constant initializers and type expressions.
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ValueSpec:
						// Determine skip names: if this is a const or var declaration,
						// skip the names being declared (self-reference is not a "use").
						var skipNames []*ast.Ident
						if d.Tok == token.CONST {
							skipNames = s.Names
						}
						for _, val := range s.Values {
							scanExprForConstReads(val, allConsts, readConsts, skipNames)
						}
						if s.Type != nil {
							scanExprForConstReads(s.Type, allConsts, readConsts, nil)
						}
					case *ast.TypeSpec:
						scanExprForConstReads(s.Type, allConsts, readConsts, nil)
					}
				}
			}
		}
	}

	// Phase 2.5: Apply iota group rule: if any constant in a block is used,
	// mark all constants in that block as used.
	for _, blockNames := range constBlocks {
		anyUsed := false
		for _, name := range blockNames {
			if readConsts[name] || entryConsts[name] {
				anyUsed = true
				break
			}
		}
		if anyUsed {
			for _, name := range blockNames {
				readConsts[name] = true
			}
		}
	}

	// Phase 3: Report unused constants declared in this file.
	thisFileName := fileNameFromFile(file)

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
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
				if entryConsts[name] {
					continue
				}

				info, exists := allConsts[name]
				if !exists {
					continue
				}
				// Only report if declared in this file.
				if info.fileName != thisFileName {
					continue
				}

				if readConsts[name] {
					continue
				}

				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryLogic,
					Confidence: 1,
					Node:       nameIdent,
					Failure:    fmt.Sprintf("const %s is unused", name),
				})
			}
		}
	}

	return failures
}

// scanForConstReads walks a statement/block looking for identifier reads that reference
// package-level constants.
func scanForConstReads(node ast.Node, allConsts map[string]*pkgConstInfo, readConsts map[string]bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			return true
		}
		if _, isPkgConst := allConsts[ident.Name]; isPkgConst {
			readConsts[ident.Name] = true
		}
		return true
	})
}

// scanExprForConstReads scans an expression for reads of package-level constants.
// skipNames is a set of ident names to skip (e.g. when scanning a const initializer,
// skip the const being declared itself).
func scanExprForConstReads(expr ast.Expr, allConsts map[string]*pkgConstInfo, readConsts map[string]bool, skipNames []*ast.Ident) {
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
		if _, isPkgConst := allConsts[name]; isPkgConst {
			readConsts[name] = true
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
func (*NoUnusedConstantRule) Name() string {
	return "noUnusedConstant"
}

// Group returns the rule group.
func (*NoUnusedConstantRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedConstantRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
