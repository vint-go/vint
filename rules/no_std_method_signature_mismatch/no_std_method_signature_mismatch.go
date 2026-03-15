package no_std_method_signature_mismatch

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// canonicalMethod describes the expected signature of a well-known interface method.
type canonicalMethod struct {
	// numParams is the number of parameters (excluding receiver).
	numParams int
	// numResults is the number of results.
	numResults int
	// resultTypes describes the expected result types (simplified string names).
	resultTypes []string
	// paramTypes describes the expected param types (simplified string names), if relevant.
	paramTypes []string
	// ifaceName is a human-readable description like "fmt.Stringer".
	ifaceName string
}

// wellKnownMethods maps method names to their expected canonical signatures.
// Based on golang.org/x/tools/go/analysis/passes/stdmethods.
var wellKnownMethods = map[string]canonicalMethod{
	// fmt.Stringer
	"String": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"string"},
		ifaceName:   "fmt.Stringer",
	},
	// fmt.GoStringer
	"GoString": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"string"},
		ifaceName:   "fmt.GoStringer",
	},
	// io.Reader
	"Read": {
		numParams:   1,
		numResults:  2,
		resultTypes: []string{"int", "error"},
		paramTypes:  []string{"[]byte"},
		ifaceName:   "io.Reader",
	},
	// io.Writer
	"Write": {
		numParams:   1,
		numResults:  2,
		resultTypes: []string{"int", "error"},
		paramTypes:  []string{"[]byte"},
		ifaceName:   "io.Writer",
	},
	// io.Closer
	"Close": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"error"},
		ifaceName:   "io.Closer",
	},
	// io.Seeker
	"Seek": {
		numParams:   2,
		numResults:  2,
		resultTypes: []string{"int64", "error"},
		paramTypes:  []string{"int64", "int"},
		ifaceName:   "io.Seeker",
	},
	// encoding/json.Marshaler
	"MarshalJSON": {
		numParams:   0,
		numResults:  2,
		resultTypes: []string{"[]byte", "error"},
		ifaceName:   "encoding/json.Marshaler",
	},
	// encoding/json.Unmarshaler
	"UnmarshalJSON": {
		numParams:   1,
		numResults:  1,
		resultTypes: []string{"error"},
		paramTypes:  []string{"[]byte"},
		ifaceName:   "encoding/json.Unmarshaler",
	},
	// encoding.TextMarshaler
	"MarshalText": {
		numParams:   0,
		numResults:  2,
		resultTypes: []string{"[]byte", "error"},
		ifaceName:   "encoding.TextMarshaler",
	},
	// encoding.TextUnmarshaler
	"UnmarshalText": {
		numParams:   1,
		numResults:  1,
		resultTypes: []string{"error"},
		paramTypes:  []string{"[]byte"},
		ifaceName:   "encoding.TextUnmarshaler",
	},
	// encoding.BinaryMarshaler
	"MarshalBinary": {
		numParams:   0,
		numResults:  2,
		resultTypes: []string{"[]byte", "error"},
		ifaceName:   "encoding.BinaryMarshaler",
	},
	// encoding.BinaryUnmarshaler
	"UnmarshalBinary": {
		numParams:   1,
		numResults:  1,
		resultTypes: []string{"error"},
		paramTypes:  []string{"[]byte"},
		ifaceName:   "encoding.BinaryUnmarshaler",
	},
	// fmt.Scanner
	"Scan": {
		numParams:   2,
		numResults:  1,
		resultTypes: []string{"error"},
		ifaceName:   "fmt.Scanner",
	},
	// fmt.Formatter
	"Format": {
		numParams:   2,
		numResults:  0,
		ifaceName:   "fmt.Formatter",
	},
	// error
	"Error": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"string"},
		ifaceName:   "error",
	},
	// io.WriterTo
	"WriteTo": {
		numParams:   1,
		numResults:  2,
		resultTypes: []string{"int64", "error"},
		ifaceName:   "io.WriterTo",
	},
	// io.ReaderFrom
	"ReadFrom": {
		numParams:   1,
		numResults:  2,
		resultTypes: []string{"int64", "error"},
		ifaceName:   "io.ReaderFrom",
	},
	// io.ReaderAt
	"ReadAt": {
		numParams:   2,
		numResults:  2,
		resultTypes: []string{"int", "error"},
		paramTypes:  []string{"[]byte", "int64"},
		ifaceName:   "io.ReaderAt",
	},
	// io.WriterAt
	"WriteAt": {
		numParams:   2,
		numResults:  2,
		resultTypes: []string{"int", "error"},
		paramTypes:  []string{"[]byte", "int64"},
		ifaceName:   "io.WriterAt",
	},
	// io.ByteReader
	"ReadByte": {
		numParams:   0,
		numResults:  2,
		resultTypes: []string{"byte", "error"},
		ifaceName:   "io.ByteReader",
	},
	// io.ByteWriter
	"WriteByte": {
		numParams:   1,
		numResults:  1,
		resultTypes: []string{"error"},
		paramTypes:  []string{"byte"},
		ifaceName:   "io.ByteWriter",
	},
	// io.ByteScanner - UnreadByte
	"UnreadByte": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"error"},
		ifaceName:   "io.ByteScanner",
	},
	// io.RuneReader
	"ReadRune": {
		numParams:   0,
		numResults:  3,
		resultTypes: []string{"rune", "int", "error"},
		ifaceName:   "io.RuneReader",
	},
	// io.RuneScanner - UnreadRune
	"UnreadRune": {
		numParams:   0,
		numResults:  1,
		resultTypes: []string{"error"},
		ifaceName:   "io.RuneScanner",
	},
	// io.StringWriter
	"WriteString": {
		numParams:   1,
		numResults:  2,
		resultTypes: []string{"int", "error"},
		paramTypes:  []string{"string"},
		ifaceName:   "io.StringWriter",
	},
}

// NoStdMethodSignatureMismatchRule checks for methods with names matching well-known interfaces
// but with incorrect signatures.
type NoStdMethodSignatureMismatchRule struct{}

// Apply applies the rule to given file.
func (r *NoStdMethodSignatureMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Only check methods (functions with a receiver).
		if funcDecl.Recv == nil {
			continue
		}

		methodName := funcDecl.Name.Name
		canonical, exists := wellKnownMethods[methodName]
		if !exists {
			continue
		}

		if !matchesCanonical(funcDecl.Type, canonical) {
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       funcDecl,
				Failure:    fmt.Sprintf("method %s has signature mismatch with %s interface", methodName, canonical.ifaceName),
			})
		}
	}

	return failures
}

// matchesCanonical checks whether a function type matches the canonical method signature.
func matchesCanonical(ft *ast.FuncType, cm canonicalMethod) bool {
	// Check parameter count.
	numParams := countFields(ft.Params)
	if numParams != cm.numParams {
		return false
	}

	// Check result count.
	numResults := countFields(ft.Results)
	if numResults != cm.numResults {
		return false
	}

	// Check result types if specified.
	if len(cm.resultTypes) > 0 && ft.Results != nil {
		resultTypes := flattenFieldTypes(ft.Results)
		if len(resultTypes) != len(cm.resultTypes) {
			return false
		}
		for i, expected := range cm.resultTypes {
			actual := typeExprToString(resultTypes[i])
			if actual != expected {
				return false
			}
		}
	}

	// Check param types if specified.
	if len(cm.paramTypes) > 0 && ft.Params != nil {
		paramTypes := flattenFieldTypes(ft.Params)
		if len(paramTypes) != len(cm.paramTypes) {
			return false
		}
		for i, expected := range cm.paramTypes {
			actual := typeExprToString(paramTypes[i])
			if actual != expected {
				return false
			}
		}
	}

	return true
}

// countFields counts the total number of individual fields (expanding grouped names).
func countFields(fl *ast.FieldList) int {
	if fl == nil {
		return 0
	}
	count := 0
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			count++
		} else {
			count += len(f.Names)
		}
	}
	return count
}

// flattenFieldTypes returns one ast.Expr per field entry (expanding grouped names).
func flattenFieldTypes(fl *ast.FieldList) []ast.Expr {
	if fl == nil {
		return nil
	}
	var types []ast.Expr
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			types = append(types, f.Type)
		} else {
			for range f.Names {
				types = append(types, f.Type)
			}
		}
	}
	return types
}

// typeExprToString converts a simple AST type expression to its string representation.
func typeExprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
		return ""
	case *ast.ArrayType:
		if t.Len == nil {
			// Slice type.
			elem := typeExprToString(t.Elt)
			return "[]" + elem
		}
		// Array type with length - render as [N]type.
		return "[...]" + typeExprToString(t.Elt)
	case *ast.StarExpr:
		return "*" + typeExprToString(t.X)
	case *ast.MapType:
		return "map[" + typeExprToString(t.Key) + "]" + typeExprToString(t.Value)
	case *ast.InterfaceType:
		if t.Methods == nil || len(t.Methods.List) == 0 {
			return "interface{}"
		}
		return "interface{...}"
	case *ast.Ellipsis:
		return "..." + typeExprToString(t.Elt)
	case *ast.ChanType:
		elem := typeExprToString(t.Value)
		switch t.Dir {
		case ast.SEND:
			return "chan<- " + elem
		case ast.RECV:
			return "<-chan " + elem
		default:
			return "chan " + elem
		}
	case *ast.FuncType:
		return "func(...)"
	case *ast.ParenExpr:
		return "(" + typeExprToString(t.X) + ")"
	case *ast.IndexExpr:
		return typeExprToString(t.X) + "[" + typeExprToString(t.Index) + "]"
	default:
		return ""
	}
}

// Name returns the rule name.
func (*NoStdMethodSignatureMismatchRule) Name() string {
	return "noStdMethodSignatureMismatch"
}

// Group returns the rule group.
func (*NoStdMethodSignatureMismatchRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoStdMethodSignatureMismatchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// canonicalSignature returns a human-readable expected signature for a method.
func canonicalSignature(name string, cm canonicalMethod) string {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString("(")
	for i, p := range cm.paramTypes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p)
	}
	b.WriteString(")")
	if cm.numResults > 0 {
		if cm.numResults == 1 {
			b.WriteString(" ")
			b.WriteString(cm.resultTypes[0])
		} else {
			b.WriteString(" (")
			for i, r := range cm.resultTypes {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(r)
			}
			b.WriteString(")")
		}
	}
	return b.String()
}
