package no_unused_field

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnusedFieldRule detects struct fields that are declared but never accessed
// (read or written to) anywhere in the package.
type NoUnusedFieldRule struct {
	fieldWritesAreUses   bool
	exportedFieldsAreUsed bool
	generatedIsUsed      bool
}

const (
	defaultFieldWritesAreUses   = true
	defaultExportedFieldsAreUsed = true
	defaultGeneratedIsUsed      = true
)

// Configure validates and applies the rule configuration.
func (r *NoUnusedFieldRule) Configure(arguments lint.Arguments) error {
	r.fieldWritesAreUses = defaultFieldWritesAreUses
	r.exportedFieldsAreUsed = defaultExportedFieldsAreUsed
	r.generatedIsUsed = defaultGeneratedIsUsed

	if len(arguments) == 0 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noUnusedField" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch normalizeOption(k) {
		case "fieldwritesareuses":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "field-writes-are-uses" in "noUnusedField"; need bool but got %T`, v)
			}
			r.fieldWritesAreUses = b
		case "exportedfieldsareused":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "exported-fields-are-used" in "noUnusedField"; need bool but got %T`, v)
			}
			r.exportedFieldsAreUsed = b
		case "generatedisused":
			b, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid value for "generated-is-used" in "noUnusedField"; need bool but got %T`, v)
			}
			r.generatedIsUsed = b
		}
	}

	return nil
}

// fieldInfo holds information about a struct field declaration.
type fieldInfo struct {
	structName string
	fieldName  string
	node       ast.Node // the *ast.Field node
	fileName   string
}

// fieldKey uniquely identifies a field within a struct.
type fieldKey struct {
	structName string
	fieldName  string
}

// Apply applies the rule to the given file.
func (r *NoUnusedFieldRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fieldWritesAreUses := r.fieldWritesAreUses
	exportedFieldsAreUsed := r.exportedFieldsAreUsed

	// Use defaults if Configure was not called.
	if !fieldWritesAreUses && !r.exportedFieldsAreUsed && !r.generatedIsUsed {
		fieldWritesAreUses = defaultFieldWritesAreUses
		exportedFieldsAreUsed = defaultExportedFieldsAreUsed
	}

	pkg := file.Pkg

	// Phase 1: Collect all struct field declarations across the package.
	allFields := map[fieldKey]*fieldInfo{}
	usedFields := map[fieldKey]bool{}
	// Track structs that have HostLayout fields -- all their fields are used.
	hostLayoutStructs := map[string]bool{}

	for fname, f := range pkg.Files() {
		// If generated-is-used and file is generated, skip collecting fields from it.
		if r.generatedIsUsed && isGeneratedFile(f) {
			continue
		}

		for _, decl := range f.AST.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				structType, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				structName := ts.Name.Name

				hasHostLayout := false
				for _, field := range structType.Fields.List {
					// Check for HostLayout sentinel (structs.HostLayout).
					if isHostLayoutField(field) {
						hasHostLayout = true
					}
				}
				if hasHostLayout {
					hostLayoutStructs[structName] = true
				}

				for _, field := range structType.Fields.List {
					// Handle embedded fields.
					if len(field.Names) == 0 {
						// Embedded fields -- skip, they are generally considered used
						// if they help implement interfaces, have exported methods/fields, etc.
						// For simplicity and to match the spec (rules 6.3, 6.4, 6.5),
						// we consider all embedded fields as used.
						continue
					}

					for _, nameIdent := range field.Names {
						fieldName := nameIdent.Name

						// Blank identifier is always considered used.
						if fieldName == "_" {
							continue
						}

						// NoCopy sentinel fields are always considered used.
						if isNoCopyField(field) {
							continue
						}

						// Exported fields are considered used by default.
						if exportedFieldsAreUsed && ast.IsExported(fieldName) {
							continue
						}

						key := fieldKey{structName: structName, fieldName: fieldName}
						allFields[key] = &fieldInfo{
							structName: structName,
							fieldName:  fieldName,
							node:       field,
							fileName:   fname,
						}
					}
				}
			}
		}
	}

	if len(allFields) == 0 {
		return nil
	}

	// Mark all fields in HostLayout structs as used.
	for key := range allFields {
		if hostLayoutStructs[key.structName] {
			usedFields[key] = true
		}
	}

	// Phase 2: Scan for field usages across the package.
	for _, f := range pkg.Files() {
		scanFileForFieldUsages(f, allFields, usedFields, fieldWritesAreUses)
	}

	// Phase 3: Report unused fields declared in this file.
	thisFileName := fileNameFromFile(file)

	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structName := ts.Name.Name

			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}
				for _, nameIdent := range field.Names {
					fieldName := nameIdent.Name
					key := fieldKey{structName: structName, fieldName: fieldName}

					info, exists := allFields[key]
					if !exists {
						continue
					}
					// Only report if declared in this file.
					if info.fileName != thisFileName {
						continue
					}

					if usedFields[key] {
						continue
					}

					failures = append(failures, lint.Failure{
						Category:   lint.FailureCategoryLogic,
						Confidence: 1,
						Node:       field,
						Failure:    fmt.Sprintf("field %s is unused", fieldName),
					})
				}
			}
		}
	}

	return failures
}

// scanFileForFieldUsages scans a file's AST for field accesses (reads and writes).
func scanFileForFieldUsages(file *lint.File, allFields map[fieldKey]*fieldInfo, usedFields map[fieldKey]bool, fieldWritesAreUses bool) {
	// Collect named struct types for struct-name resolution from composite literals.
	structTypes := map[string]*ast.StructType{}
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if st, ok := ts.Type.(*ast.StructType); ok {
				structTypes[ts.Name.Name] = st
			}
		}
	}

	ast.Inspect(file.AST, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Field access: expr.FieldName
			fieldName := x.Sel.Name

			// Try to determine the struct type from the receiver expression.
			structName := resolveStructName(x.X, structTypes, file)

			if structName != "" {
				key := fieldKey{structName: structName, fieldName: fieldName}
				if _, tracked := allFields[key]; tracked {
					// Check if this is a write-only access.
					if fieldWritesAreUses || !isWriteOnlyAccess(x, file) {
						usedFields[key] = true
					}
				}
			} else {
				// If we cannot determine the struct type, mark any field
				// with this name as used (conservative approach).
				for key := range allFields {
					if key.fieldName == fieldName {
						if fieldWritesAreUses || !isWriteOnlyAccess(x, file) {
							usedFields[key] = true
						}
					}
				}
			}

		case *ast.CompositeLit:
			// Composite literal: SomeStruct{field: value}
			structName := resolveCompositeLitStructName(x)
			if structName == "" {
				return true
			}

			// Check for keyed field initializations.
			for _, elt := range x.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					// Positional field init -- all fields at those positions are used.
					markAllFieldsOfStructUsed(structName, allFields, usedFields)
					break
				}
				ident, ok := kv.Key.(*ast.Ident)
				if !ok {
					continue
				}
				key := fieldKey{structName: structName, fieldName: ident.Name}
				if _, tracked := allFields[key]; tracked {
					if fieldWritesAreUses {
						usedFields[key] = true
					}
				}
			}

		case *ast.CallExpr:
			// Check for unsafe.Pointer conversions -- all fields become used.
			if isUnsafePointerConversion(x) {
				markAllFieldsUsedForExpr(x, structTypes, allFields, usedFields, file)
			}
		}
		return true
	})
}

// resolveStructName attempts to determine the struct type name from an expression
// used as the receiver of a field access (e.g., the "x" in "x.Field").
func resolveStructName(expr ast.Expr, structTypes map[string]*ast.StructType, file *lint.File) string {
	switch x := expr.(type) {
	case *ast.Ident:
		// Direct variable access: look up the variable's type from its declaration.
		if x.Obj != nil {
			switch decl := x.Obj.Decl.(type) {
			case *ast.AssignStmt:
				// Short variable declaration: v := SomeStruct{}
				for i, lhs := range decl.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name == x.Name {
						if i < len(decl.Rhs) {
							return resolveExprStructName(decl.Rhs[i], structTypes)
						}
					}
				}
			case *ast.ValueSpec:
				// var v SomeStruct
				if decl.Type != nil {
					return resolveTypeExprName(decl.Type)
				}
				// var v = SomeStruct{}
				for _, val := range decl.Values {
					name := resolveExprStructName(val, structTypes)
					if name != "" {
						return name
					}
				}
			case *ast.Field:
				// Function parameter: func foo(v SomeStruct)
				return resolveTypeExprName(decl.Type)
			}
		}
	case *ast.UnaryExpr:
		if x.Op == token.AND {
			return resolveStructName(x.X, structTypes, file)
		}
	case *ast.StarExpr:
		return resolveStructName(x.X, structTypes, file)
	case *ast.CompositeLit:
		return resolveCompositeLitStructName(x)
	}
	return ""
}

// resolveExprStructName tries to determine the struct type name from a value expression.
func resolveExprStructName(expr ast.Expr, structTypes map[string]*ast.StructType) string {
	switch x := expr.(type) {
	case *ast.CompositeLit:
		return resolveCompositeLitStructName(x)
	case *ast.UnaryExpr:
		if x.Op == token.AND {
			return resolveExprStructName(x.X, structTypes)
		}
	}
	return ""
}

// resolveCompositeLitStructName extracts the struct type name from a composite literal.
func resolveCompositeLitStructName(lit *ast.CompositeLit) string {
	if lit.Type == nil {
		return ""
	}
	return resolveTypeExprName(lit.Type)
}

// resolveTypeExprName extracts the type name from a type expression.
func resolveTypeExprName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return resolveTypeExprName(t.X)
	case *ast.SelectorExpr:
		// pkg.Type -- external type, not our concern.
		return ""
	}
	return ""
}

// isWriteOnlyAccess checks if the given selector expression is used only
// as the target of an assignment (write-only), not as a value (read).
func isWriteOnlyAccess(sel *ast.SelectorExpr, file *lint.File) bool {
	// To keep this simple and avoid false positives, we conservatively
	// return false (treat as read) in most cases.
	// A full implementation would walk the parent nodes.
	return false
}

// markAllFieldsOfStructUsed marks all tracked fields of a struct as used.
func markAllFieldsOfStructUsed(structName string, allFields map[fieldKey]*fieldInfo, usedFields map[fieldKey]bool) {
	for key := range allFields {
		if key.structName == structName {
			usedFields[key] = true
		}
	}
}

// markAllFieldsUsedForExpr marks all tracked fields as used for expressions
// involved in unsafe.Pointer conversions.
func markAllFieldsUsedForExpr(call *ast.CallExpr, structTypes map[string]*ast.StructType, allFields map[fieldKey]*fieldInfo, usedFields map[fieldKey]bool, file *lint.File) {
	for _, arg := range call.Args {
		structName := resolveStructName(arg, structTypes, file)
		if structName != "" {
			markAllFieldsOfStructUsed(structName, allFields, usedFields)
		}
	}
	// Also mark all fields used conservatively.
	for key := range allFields {
		usedFields[key] = true
	}
}

// isUnsafePointerConversion checks if a call expression is a conversion
// to or from unsafe.Pointer.
func isUnsafePointerConversion(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if ident.Name == "unsafe" && sel.Sel.Name == "Pointer" {
				return true
			}
		}
	}
	return false
}

// isNoCopyField checks if a field's type is a NoCopy sentinel type.
func isNoCopyField(field *ast.Field) bool {
	typeName := resolveTypeExprName(field.Type)
	// Common NoCopy sentinel types.
	return typeName == "noCopy" || typeName == "NoCopy"
}

// isHostLayoutField checks if a field is a structs.HostLayout sentinel.
func isHostLayoutField(field *ast.Field) bool {
	// Check for embedded structs.HostLayout field.
	if len(field.Names) == 0 {
		sel, ok := field.Type.(*ast.SelectorExpr)
		if ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				return ident.Name == "structs" && sel.Sel.Name == "HostLayout"
			}
		}
	}
	return false
}

// isGeneratedFile checks if a file has a "Code generated" comment.
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

// fileNameFromFile extracts the file name that matches pkg.Files() keys.
func fileNameFromFile(file *lint.File) string {
	for name, f := range file.Pkg.Files() {
		if f == file {
			return name
		}
	}
	return ""
}

// normalizeOption normalizes a configuration option name by removing hyphens,
// underscores, and lowering case.
func normalizeOption(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return strings.ToLower(name)
}

// Name returns the rule name.
func (*NoUnusedFieldRule) Name() string {
	return "noUnusedField"
}

// Group returns the rule group.
func (*NoUnusedFieldRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnusedFieldRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}
