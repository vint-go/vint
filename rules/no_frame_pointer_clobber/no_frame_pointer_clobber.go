package no_frame_pointer_clobber

import (
	"bufio"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoFramePointerClobberRule reports assembly code that clobbers the frame pointer
// before saving it. On architectures that use a frame pointer (such as amd64),
// assembly functions must save the frame pointer (BP register) before modifying it.
// Failing to do so breaks stack unwinding, which affects debugging, profiling,
// and panic stack traces.
type NoFramePointerClobberRule struct{}

// asmFramePointerInfo holds information about a frame pointer clobber in assembly.
type asmFramePointerInfo struct {
	funcName string
	line     int
	message  string
}

// textDirectiveRegexp matches assembly TEXT directives for frameless functions.
// A frameless function has frame size $0, which means it does not allocate
// stack space and therefore must save BP before modifying it.
// Examples:
//
//	TEXT ·example(SB), NOSPLIT, $0-8
//	TEXT ·example(SB), $0
//	TEXT ·example(SB),NOSPLIT,$0-8
var textDirectiveRegexp = regexp.MustCompile(
	`(?m)^\s*TEXT\s+` + // TEXT keyword
		`(?:\w+)?` + // optional package name
		"\xC2\xB7" + // middle dot separator (U+00B7 in UTF-8)
		`(\w+)` + // function name (capture group 1)
		`\(SB\)` + // (SB) suffix
		`\s*,\s*` + // comma separator
		`(?:[^,]*,\s*)?` + // optional flags like NOSPLIT
		`\$0` + // frame size must be $0 (frameless)
		`(?:-\d+)?`, // optional arg size
)

// fpWriteRegexp matches instructions that write to the BP register (amd64).
// Matches lines where BP appears as a destination operand (at the end of instruction).
var fpWriteRegexp = regexp.MustCompile(`,\s*BP$`)

// fpReadRegexp matches instructions that read the BP register.
var fpReadRegexp = regexp.MustCompile(`\bBP\b`)

// unconditionalBranchRegexp matches unconditional branch instructions.
var unconditionalBranchRegexp = regexp.MustCompile(`^\s*(JMP|RET)\b`)

// Apply applies the rule to the given file.
func (r *NoFramePointerClobberRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Find Go function declarations without bodies (assembly stubs).
	stubFuncs := map[string]*ast.FuncDecl{}
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		// Functions without a body are assembly stubs.
		if funcDecl.Body == nil && funcDecl.Recv == nil {
			stubFuncs[funcDecl.Name.Name] = funcDecl
		}
	}

	if len(stubFuncs) == 0 {
		return nil
	}

	// Look for assembly files in the same directory and check for frame pointer clobbers.
	dir := filepath.Dir(file.Name)
	clobbers := findFramePointerClobbers(dir)
	if len(clobbers) == 0 {
		return nil
	}

	// Report failures on corresponding Go stub declarations.
	for _, clobber := range clobbers {
		funcDecl, found := stubFuncs[clobber.funcName]
		if !found {
			continue
		}
		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       funcDecl,
			Failure:    clobber.message,
		})
	}

	return failures
}

// findFramePointerClobbers scans assembly files in the given directory for
// frameless functions that clobber BP before saving it.
func findFramePointerClobbers(dir string) []asmFramePointerInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var result []asmFramePointerInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".s") {
			continue
		}
		clobbers := checkAsmFileForClobbers(filepath.Join(dir, entry.Name()))
		result = append(result, clobbers...)
	}
	return result
}

// checkAsmFileForClobbers checks a single assembly file for frame pointer clobber issues.
// For each frameless TEXT directive, it checks whether BP is written before being read
// or before an unconditional branch.
func checkAsmFileForClobbers(path string) []asmFramePointerInfo {
	f, err := os.Open(path) //nolint:gosec // we need to read assembly files
	if err != nil {
		return nil
	}
	defer f.Close()

	var result []asmFramePointerInfo
	scanner := bufio.NewScanner(f)

	// State: when active is true, we are inside a frameless function
	// and watching for BP clobber before save.
	active := false
	currentFunc := ""
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Skip comments and empty lines.
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		// Check for a new TEXT directive.
		matches := textDirectiveRegexp.FindStringSubmatch(line)
		if matches != nil {
			funcName := matches[1]
			// Start watching for BP clobber in this frameless function.
			active = true
			currentFunc = funcName
			continue
		}

		if !active {
			continue
		}

		// Check if this is a new TEXT directive (different function start).
		if strings.HasPrefix(trimmed, "TEXT") {
			active = false
			continue
		}

		// If BP is written (as a destination), it's clobbered.
		if fpWriteRegexp.MatchString(trimmed) {
			result = append(result, asmFramePointerInfo{
				funcName: currentFunc,
				line:     lineNum,
				message: fmt.Sprintf(
					"assembly function %s clobbers frame pointer (BP) before saving it",
					currentFunc),
			})
			active = false
			continue
		}

		// If BP is read, it means it's being saved/used properly.
		if fpReadRegexp.MatchString(trimmed) {
			active = false
			continue
		}

		// If we hit an unconditional branch, stop checking this function.
		if unconditionalBranchRegexp.MatchString(trimmed) {
			active = false
			continue
		}
	}

	return result
}

// Name returns the rule name.
func (*NoFramePointerClobberRule) Name() string {
	return "noFramePointerClobber"
}

// Group returns the rule group.
func (*NoFramePointerClobberRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
// This rule reads assembly files from disk alongside the Go file,
// so it is not purely file-only. However, it doesn't use type checking.
// We use TierFileOnly because assembly files are external to the package
// graph and the rule just needs the Go file's AST.
func (*NoFramePointerClobberRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
