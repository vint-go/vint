package lint

import (
	"bytes"
	"fmt"
	"go/build"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/sync/singleflight"
	"golang.org/x/tools/go/gcexportdata"
)

// safeImporter is a goroutine-safe alternative to sharedImporter.
// It uses gcexportdata.Read under a mutex to avoid the concurrent map
// access race in go/internal/gcimporter (see bug-goimporter.md).
//
// File finding + I/O runs concurrently (no shared mutable state).
// Only the fast CPU-bound gcexportdata.Read is serialized under mu.
type safeImporter struct {
	cache sync.Map           // import path → *importResult (lock-free reads)
	sf    singleflight.Group // deduplicates concurrent resolution of the same import path
	fset  *token.FileSet

	mu       sync.Mutex                // protects packages map during gcexportdata.Read
	packages map[string]*types.Package // shared across all Read calls

	exportCache sync.Map // import path → export file path (string), caches findExportFile results
}

func newSafeImporter() *safeImporter {
	return &safeImporter{
		fset:     token.NewFileSet(),
		packages: make(map[string]*types.Package),
	}
}

func (s *safeImporter) Import(path string) (*types.Package, error) {
	return s.ImportFrom(path, "", 0)
}

func (s *safeImporter) ImportFrom(path, srcDir string, mode types.ImportMode) (*types.Package, error) {
	// Fast path: lock-free cache check.
	if v, ok := s.cache.Load(path); ok {
		r := v.(*importResult)
		return r.pkg, r.err
	}

	// The "unsafe" package is built-in; it has no export data file.
	if path == "unsafe" {
		return types.Unsafe, nil
	}

	v, err, _ := s.sf.Do(path, func() (interface{}, error) {
		// Double-check cache after entering singleflight.
		if v, ok := s.cache.Load(path); ok {
			return v, nil
		}

		// Step 1: Find and read export data (concurrent, no shared state).
		data, readErr := s.findAndReadExport(path, srcDir)
		if readErr != nil {
			r := &importResult{err: readErr}
			s.cache.Store(path, r)
			return r, nil
		}

		// Step 2: Parse export data under mutex (fast, CPU-bound).
		s.mu.Lock()
		// Another singleflight call for a different path may have
		// populated this package as a transitive dependency.
		if pkg := s.packages[path]; pkg != nil && pkg.Complete() {
			s.mu.Unlock()
			r := &importResult{pkg: pkg}
			s.cache.Store(path, r)
			return r, nil
		}
		pkg, parseErr := gcexportdata.Read(bytes.NewReader(data), s.fset, s.packages, path)
		s.mu.Unlock()

		r := &importResult{pkg: pkg, err: parseErr}
		s.cache.Store(path, r)
		return r, nil
	})
	if err != nil {
		return nil, err
	}

	r := v.(*importResult)
	return r.pkg, r.err
}

// findAndReadExport locates the export data file for path and reads
// the raw export data bytes. This runs concurrently without any locks.
func (s *safeImporter) findAndReadExport(path, srcDir string) ([]byte, error) {
	filename, err := s.findExportFile(path, srcDir)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := gcexportdata.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("reading export data for %q: %v", path, err)
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading export data for %q: %v", path, err)
	}
	return data, nil
}

// findExportFile locates the .a archive containing export data for
// the given import path. It replicates the logic of the stdlib's
// internal/exportdata.FindPkg using public APIs.
func (s *safeImporter) findExportFile(path, srcDir string) (string, error) {
	// Check cache first.
	if v, ok := s.exportCache.Load(path); ok {
		return v.(string), nil
	}

	bp, err := build.Import(path, srcDir, build.FindOnly|build.AllowBinary)
	if err != nil {
		return "", fmt.Errorf("can't find import: %s: %w", path, err)
	}

	// Try the installed .a file path.
	if bp.PkgObj != "" {
		if _, statErr := os.Stat(bp.PkgObj); statErr == nil {
			s.exportCache.Store(path, bp.PkgObj)
			return bp.PkgObj, nil
		}
	}

	// For standard library packages, the .a file may be in the build
	// cache rather than at PkgObj. Use "go list -export" to find it.
	if bp.Goroot {
		filename, err := lookupGorootExport(path, srcDir)
		if err != nil {
			return "", err
		}
		s.exportCache.Store(path, filename)
		return filename, nil
	}

	return "", fmt.Errorf("no export data for %q", path)
}

// lookupGorootExport uses "go list -export" to find the export data
// file for a standard library package. This spawns a subprocess but
// runs concurrently (outside the main mutex).
func lookupGorootExport(importPath, srcDir string) (string, error) {
	cmd := exec.Command("go", "list", "-export", "-f", "{{.Export}}", "--", importPath)
	if srcDir != "" {
		cmd.Dir = srcDir
	}
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go list -export %s: %w", importPath, err)
	}
	result := strings.TrimSpace(string(out))
	if result == "" {
		return "", fmt.Errorf("no export file for %q", importPath)
	}
	return result, nil
}

// cachedResult returns the cached import result for the given path.
// Used by ImportedPkgSourceDir.
func (s *safeImporter) cachedResult(path string) (*importResult, bool) {
	v, ok := s.cache.Load(path)
	if !ok {
		return nil, false
	}
	return v.(*importResult), true
}

// importerFileSet returns the FileSet used by this importer.
// Used by ImportedPkgSourceDir.
func (s *safeImporter) importerFileSet() *token.FileSet {
	return s.fset
}
