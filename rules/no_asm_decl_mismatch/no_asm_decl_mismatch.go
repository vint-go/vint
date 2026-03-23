package no_asm_decl_mismatch

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoAsmDeclMismatchRule reports mismatches between assembly (.s) files and Go declarations.
// It checks that assembly functions match their Go function prototypes, ensuring that
// parameter sizes, offsets, and return values are consistent.
type NoAsmDeclMismatchRule struct{}

// asmFunc holds parsed assembly function information from a TEXT directive.
type asmFunc struct {
	name    string // function name (e.g. "Add")
	argSize int    // argument size from TEXT directive (e.g. 24 in $0-24)
	hasSize bool   // whether argSize was parsed
}

// textDirectiveRegexp matches assembly TEXT directives.
// Examples:
//
//	TEXT ·Add(SB), NOSPLIT, $0-24
//	TEXT ·Add(SB), $0-24
//	TEXT ·Add(SB),NOSPLIT,$0-24
var textDirectiveRegexp = regexp.MustCompile(
	`(?m)^\s*TEXT\s+` + // TEXT keyword
		`(?:\w+)?` + // optional package name
		"\xC2\xB7" + // middle dot separator (U+00B7 in UTF-8)
		`(\w+)` + // function name (capture group 1)
		`\(SB\)` + // (SB) suffix
		`\s*,\s*` + // comma separator
		`(?:[^,]*,\s*)?` + // optional flags like NOSPLIT
		`\$(\d+)` + // frame size (capture group 2)
		`(?:-(\d+))?`, // optional arg size (capture group 3)
)

// Apply applies the rule to given file.
func (r *NoAsmDeclMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Find Go function declarations without bodies (assembly stubs).
	var stubFuncs []*ast.FuncDecl
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Functions without a body are assembly stubs.
		if funcDecl.Body == nil && funcDecl.Recv == nil {
			stubFuncs = append(stubFuncs, funcDecl)
		}
	}

	if len(stubFuncs) == 0 {
		return nil
	}

	// Look for assembly files in the same directory.
	dir := filepath.Dir(file.Name)
	asmFuncs := findAsmFunctions(dir)
	if len(asmFuncs) == 0 {
		return nil
	}

	// Enable type checking for computing parameter sizes.
	file.Pkg.TypeCheck() //nolint:errcheck // partial type info is fine

	typesInfo := file.Pkg.TypesInfo()

	// Check each stub function against assembly definitions.
	for _, funcDecl := range stubFuncs {
		funcName := funcDecl.Name.Name
		asmDef, found := asmFuncs[funcName]
		if !found || !asmDef.hasSize {
			continue
		}

		// Compute expected argument size from the Go function signature.
		expectedSize := computeArgSize(funcDecl, typesInfo)
		if expectedSize < 0 {
			continue // couldn't determine size
		}

		if asmDef.argSize != expectedSize {
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryLogic,
				Confidence: 1,
				Node:       funcDecl,
				Failure: fmt.Sprintf(
					"assembly function %s has argument size %d, but Go declaration expects %d",
					funcName, asmDef.argSize, expectedSize),
			})
		}
	}

	return failures
}

// findAsmFunctions scans assembly files in the given directory for TEXT directives.
func findAsmFunctions(dir string) map[string]asmFunc {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	result := map[string]asmFunc{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".s") {
			continue
		}
		parseAsmFile(filepath.Join(dir, entry.Name()), result)
	}
	return result
}

// parseAsmFile parses a single assembly file for TEXT directives.
func parseAsmFile(path string, result map[string]asmFunc) {
	f, err := os.Open(path) //nolint:gosec // we need to read assembly files
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "TEXT") {
			continue
		}
		matches := textDirectiveRegexp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		funcName := matches[1]
		af := asmFunc{name: funcName}

		if matches[3] != "" {
			argSize, err := strconv.Atoi(matches[3])
			if err == nil {
				af.argSize = argSize
				af.hasSize = true
			}
		}

		result[funcName] = af
	}
}

// computeArgSize calculates the total size of function arguments and return values
// in bytes, matching what the assembly TEXT directive expects.
func computeArgSize(funcDecl *ast.FuncDecl, info *types.Info) int {
	if info == nil {
		return -1
	}

	obj := info.Defs[funcDecl.Name]
	if obj == nil {
		return -1
	}

	fn, ok := obj.(*types.Func)
	if !ok {
		return -1
	}

	sig := fn.Type().(*types.Signature)
	total := 0

	// Sum parameter sizes.
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		size := typeSize(params.At(i).Type())
		if size < 0 {
			return -1
		}
		total += size
	}

	// Sum return value sizes.
	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		size := typeSize(results.At(i).Type())
		if size < 0 {
			return -1
		}
		total += size
	}

	return total
}

// typeSize returns the size in bytes of a type for the purposes of assembly
// argument sizing. Uses 8-byte alignment assumptions for 64-bit architectures.
func typeSize(t types.Type) int {
	switch typ := t.Underlying().(type) {
	case *types.Basic:
		return basicTypeSize(typ)
	case *types.Pointer:
		return 8
	case *types.Slice:
		return 24 // pointer + len + cap
	case *types.Struct:
		total := 0
		for i := 0; i < typ.NumFields(); i++ {
			size := typeSize(typ.Field(i).Type())
			if size < 0 {
				return -1
			}
			total += size
		}
		return total
	case *types.Array:
		elemSize := typeSize(typ.Elem())
		if elemSize < 0 {
			return -1
		}
		return elemSize * int(typ.Len())
	case *types.Interface:
		return 16 // two words: type pointer + data pointer
	case *types.Map:
		return 8 // pointer to runtime map
	case *types.Chan:
		return 8 // pointer to runtime channel
	case *types.Signature:
		return 8 // function pointer
	default:
		return -1
	}
}

// basicTypeSize returns the size of a basic type.
func basicTypeSize(t *types.Basic) int {
	switch t.Kind() {
	case types.Bool, types.Int8, types.Uint8:
		return 1
	case types.Int16, types.Uint16:
		return 2
	case types.Int32, types.Uint32, types.Float32:
		return 4
	case types.Int64, types.Uint64, types.Float64, types.Complex64:
		return 8
	case types.Int, types.Uint, types.Uintptr:
		return 8 // assume 64-bit
	case types.Complex128:
		return 16
	case types.String:
		return 16 // pointer + length
	case types.UnsafePointer:
		return 8
	default:
		return -1
	}
}

// Name returns the rule name.
func (*NoAsmDeclMismatchRule) Name() string {
	return "noAsmDeclMismatch"
}

// Group returns the rule group.
func (*NoAsmDeclMismatchRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
// This rule uses type information (TypeCheck) so it needs TierPackageAware.
func (*NoAsmDeclMismatchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoAsmDeclMismatchRule) RequiresTypecheck() bool {
	return true
}
