// import_filter_analysis scans Go files in a directory and reports
// how many AST walks could be skipped by import-based rule filtering.
//
// Usage: go run ./tools/import_filter_analysis <directory>
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ruleGroup maps a required import path to the rules that only need to
// run when that import is present. Built from inspecting rule source code.
type ruleGroup struct {
	importPath string
	rules      []string
}

var ruleGroups = []ruleGroup{
	// --- net/http family (largest group) ---
	{"net/http", []string{
		"noHttpGet", "noHttpHead", "noHttpPost", "noHttpPostForm",
		"noHttpNewRequest",
		"noHttpClientGet", "noHttpClientHead", "noHttpClientPost", "noHttpClientPostForm",
		"noHttpResponseMisuse", "noUnclosedBodies",
		"noConflictingHttpMuxPatterns", "noUnboundedFormParsing",
		"noUnsafeCorsBypass", "noUnsafeRedirectPolicy",
		"noFilesystemRootServing", "noInsecureCookie", "noServeWithoutTimeout",
		"noSsrfViaVariable", "noXssTaint",
	}},
	{"net/http/httptest", []string{
		"noHttptestNewRequest",
	}},

	// --- database/sql family ---
	{"database/sql", []string{
		"noSqlDbBegin", "noSqlDbExec", "noSqlDbPing", "noSqlDbPrepare",
		"noSqlDbQuery", "noSqlDbQueryRow",
		"noSqlStmtExec", "noSqlStmtQuery", "noSqlStmtQueryRow",
		"noSqlTxExec", "noSqlTxPrepare", "noSqlTxQuery", "noSqlTxQueryRow", "noSqlTxStmt",
	}},

	// --- net family ---
	{"net", []string{
		"noNetDial", "noNetDialTimeout",
		"noNetListen", "noNetListenPacket",
		"noNetLookupAddr", "noNetLookupCname", "noNetLookupHost", "noNetLookupIp",
		"noNetLookupMx", "noNetLookupNs", "noNetLookupPort", "noNetLookupSrv", "noNetLookupTxt",
		"noBindToAllInterfaces", "useJoinHostPort",
	}},

	// --- crypto/tls ---
	{"crypto/tls", []string{
		"noTlsConnHandshake", "noTlsDial", "noTlsDialWithDialer",
		"noTlsSessionResumptionBypass", "noInsecureTlsConfig",
	}},

	// --- sync ---
	{"sync", []string{
		"noCopiedLock", "noWaitGroupMisuse", "noInlineSyncOnceFunc",
		"waitgroupByValue",
	}},

	// --- sync/atomic ---
	{"sync/atomic", []string{
		"noAtomicAlignmentIssue", "noAtomicAssignMisuse",
	}},

	// --- reflect ---
	{"reflect", []string{
		"noDeepEqualErrors", "noReflectValueCompare",
	}},

	// --- os/exec ---
	{"os/exec", []string{
		"noExecCommand", "noVariableCommandExecution", "noCommandInjectionTaint",
	}},

	// --- os/signal ---
	{"os/signal", []string{
		"noUnbufferedSignalChannel",
	}},

	// --- sort ---
	{"sort", []string{
		"noBadSortUsage", "noInvalidSortSliceArg",
	}},

	// --- errors ---
	{"errors", []string{
		"noInvalidErrorsAs", "noDirectErrorComparison",
	}},

	// --- context ---
	{"context", []string{
		"noLostCancel",
		// context_as_argument and context_keys_type also, but they're more general
	}},

	// --- log/slog ---
	{"log/slog", []string{
		"noSlogKeyValueMismatch",
	}},

	// --- time ---
	{"time", []string{
		"noDeferTimeMisuse", "noIncorrectTimeFormat",
	}},

	// --- crypto weak hash ---
	{"crypto/md5", []string{"noWeakCryptoHash_md5"}},
	{"crypto/sha1", []string{"noWeakCryptoHash_sha1"}},

	// --- crypto weak encryption ---
	{"crypto/des", []string{"noWeakEncryptionAlgorithm_des"}},
	{"crypto/rc4", []string{"noWeakEncryptionAlgorithm_rc4"}},

	// --- templates ---
	{"html/template", []string{"noUnescapedHtmlTemplate", "noTemplateInjection_html"}},
	{"text/template", []string{"noTemplateInjection_text"}},

	// --- encoding ---
	{"encoding/json", []string{"noNonPointerUnmarshal_json", "noUnsafeDeserialization_json"}},
	{"encoding/xml", []string{"noNonPointerUnmarshal_xml", "noUnsafeDeserialization_xml"}},
	{"encoding/gob", []string{"noUnsafeDeserialization_gob"}},

	// --- archive ---
	{"archive/zip", []string{"noZipSlip_zip"}},
	{"archive/tar", []string{"noZipSlip_tar"}},

	// --- testing ---
	{"testing", []string{
		"noTestFatalInGoroutine",
	}},

	// --- fmt (for format-checking rules) ---
	{"fmt", []string{
		"noPrintfFormatMismatch",
		"noSqlFormatString",
	}},

	// --- unsafe ---
	{"unsafe", []string{
		"noInvalidUnsafePointer",
	}},

	// --- regexp ---
	{"regexp", []string{
		"noBadRegexpPattern",
	}},

	// --- strconv ---
	{"strconv", []string{
		"noStringIntConversion",
	}},

	// --- cgo ---
	{"C", []string{
		"noCgoPointerViolation",
	}},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <directory>\n", os.Args[0])
		os.Exit(1)
	}
	dir := os.Args[1]

	// Collect all .go files.
	var goFiles []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		// Skip vendor, testdata, hidden dirs.
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || strings.HasPrefix(base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})

	fmt.Printf("Scanning %d .go files in %s\n\n", len(goFiles), dir)

	// Parse imports from all files.
	fset := token.NewFileSet()
	type fileImports struct {
		path    string
		imports map[string]bool
	}
	var files []fileImports

	for _, path := range goFiles {
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			continue
		}
		imports := make(map[string]bool)
		for _, imp := range f.Imports {
			// Trim quotes from import path.
			p := strings.Trim(imp.Path.Value, `"`)
			imports[p] = true
		}
		files = append(files, fileImports{path: path, imports: imports})
	}

	fmt.Printf("Successfully parsed %d files\n\n", len(files))

	totalFiles := len(files)

	// For each rule group, count files that import the package.
	type result struct {
		importPath   string
		rules        []string
		filesWith    int
		filesWithout int
		walksSkipped int // filesWithout * len(rules)
	}
	var results []result

	totalWalksBaseline := 0
	totalWalksSkipped := 0

	for _, rg := range ruleGroups {
		filesWith := 0
		for _, f := range files {
			if f.imports[rg.importPath] {
				filesWith++
			}
		}
		filesWithout := totalFiles - filesWith
		skipped := filesWithout * len(rg.rules)
		baseline := totalFiles * len(rg.rules)
		totalWalksBaseline += baseline
		totalWalksSkipped += skipped

		results = append(results, result{
			importPath:   rg.importPath,
			rules:        rg.rules,
			filesWith:    filesWith,
			filesWithout: filesWithout,
			walksSkipped: skipped,
		})
	}

	// Sort by walks skipped (descending).
	sort.Slice(results, func(i, j int) bool {
		return results[i].walksSkipped > results[j].walksSkipped
	})

	// Print results.
	fmt.Printf("%-25s %6s %6s %6s %8s %6s\n",
		"IMPORT PATH", "RULES", "WITH", "W/OUT", "SKIPPED", "SKIP%")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range results {
		skipPct := float64(r.walksSkipped) / float64(totalFiles*len(r.rules)) * 100
		fmt.Printf("%-25s %6d %6d %6d %8d %5.1f%%\n",
			r.importPath,
			len(r.rules),
			r.filesWith,
			r.filesWithout,
			r.walksSkipped,
			skipPct,
		)
	}

	fmt.Println(strings.Repeat("-", 80))

	totalRulesInGroups := 0
	for _, rg := range ruleGroups {
		totalRulesInGroups += len(rg.rules)
	}

	// Approximate total rule count (from rules/ + rule/ directories).
	// Update this if you add/remove rules.
	const estimatedTotalRules = 319
	generalRules := estimatedTotalRules - totalRulesInGroups
	generalWalks := generalRules * totalFiles
	allWalksBaseline := estimatedTotalRules * totalFiles
	allWalksAfter := allWalksBaseline - totalWalksSkipped

	fmt.Printf("\nImport-filterable rules only:\n")
	fmt.Printf("  Filterable rules:            %d\n", totalRulesInGroups)
	fmt.Printf("  Baseline walks (no filter):  %d\n", totalWalksBaseline)
	fmt.Printf("  Walks that can be skipped:   %d\n", totalWalksSkipped)
	fmt.Printf("  Skip ratio:                  %.1f%%\n",
		float64(totalWalksSkipped)/float64(totalWalksBaseline)*100)

	fmt.Printf("\nAll rules (estimated %d total):\n", estimatedTotalRules)
	fmt.Printf("  General rules (always run):  %d\n", generalRules)
	fmt.Printf("  General walks (always):      %d\n", generalWalks)
	fmt.Printf("  All walks baseline:          %d\n", allWalksBaseline)
	fmt.Printf("  All walks after filtering:   %d\n", allWalksAfter)
	fmt.Printf("  Overall skip ratio:          %.1f%%\n",
		float64(totalWalksSkipped)/float64(allWalksBaseline)*100)
	fmt.Printf("  Overall reduction:           %d → %d walks (%.1fx less)\n",
		allWalksBaseline, allWalksAfter,
		float64(allWalksBaseline)/float64(allWalksAfter))
}
