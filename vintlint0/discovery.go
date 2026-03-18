package vintlint0

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// DiscoverPackages resolves include/exclude patterns into groups of .go files
// per directory (same [][]string format as dots.ResolvePackages).
//
// Pattern semantics:
//   - "./..." or "path/..."  : recursive walk from prefix directory
//   - "." or "path"          : single directory listing
//
// Exclude patterns follow the same semantics. Matching is by directory prefix.
func DiscoverPackages(includes, excludes []string) ([][]string, error) {
	if len(includes) == 0 {
		includes = []string{"."}
	}

	// Normalize exclude patterns to absolute directory prefixes.
	excludeDirs := make([]string, 0, len(excludes))
	for _, ex := range excludes {
		dir := strings.TrimSuffix(ex, "/...")
		dir = strings.TrimSuffix(dir, "...")
		dir = strings.TrimSuffix(dir, "/")
		dir = strings.TrimSuffix(dir, "\\")
		if dir == "" {
			dir = "."
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, fmt.Errorf("resolving exclude %q: %w", ex, err)
		}
		excludeDirs = append(excludeDirs, abs)
	}

	// dirFiles accumulates .go files grouped by directory.
	// Map key is absolute directory path.
	dirFiles := map[string][]string{}

	for _, pattern := range includes {
		recursive := strings.HasSuffix(pattern, "/...") ||
			strings.HasSuffix(pattern, "\\...") ||
			pattern == "..."

		prefix := pattern
		prefix = strings.TrimSuffix(prefix, "/...")
		prefix = strings.TrimSuffix(prefix, "\\...")
		prefix = strings.TrimSuffix(prefix, "...")
		if prefix == "" {
			prefix = "."
		}

		absPrefix, err := filepath.Abs(prefix)
		if err != nil {
			return nil, fmt.Errorf("resolving include %q: %w", pattern, err)
		}

		if recursive {
			err = filepath.WalkDir(absPrefix, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					// Skip excluded directories entirely.
					if isExcludedDir(path, excludeDirs) {
						return fs.SkipDir
					}
					// Skip hidden directories (e.g., .git).
					if name := d.Name(); name != "." && strings.HasPrefix(name, ".") {
						return fs.SkipDir
					}
					// Skip testdata directories.
					if d.Name() == "testdata" {
						return fs.SkipDir
					}
					return nil
				}
				if isGoFile(d.Name()) {
					dir := filepath.Dir(path)
					dirFiles[dir] = append(dirFiles[dir], path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("walking %q: %w", absPrefix, err)
			}
		} else {
			// Single directory listing.
			if isExcludedDir(absPrefix, excludeDirs) {
				continue
			}
			entries, err := os.ReadDir(absPrefix)
			if err != nil {
				return nil, fmt.Errorf("reading directory %q: %w", absPrefix, err)
			}
			for _, e := range entries {
				if !e.IsDir() && isGoFile(e.Name()) {
					fullPath := filepath.Join(absPrefix, e.Name())
					dirFiles[absPrefix] = append(dirFiles[absPrefix], fullPath)
				}
			}
		}
	}

	// Get working directory for making paths relative.
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}

	// Convert map to sorted slice-of-slices with relative paths.
	dirs := make([]string, 0, len(dirFiles))
	for dir := range dirFiles {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	result := make([][]string, 0, len(dirs))
	for _, dir := range dirs {
		files := dirFiles[dir]
		sort.Strings(files)
		// Convert absolute paths to relative paths.
		relFiles := make([]string, len(files))
		for i, f := range files {
			rel, err := filepath.Rel(cwd, f)
			if err != nil {
				relFiles[i] = f // fallback to absolute
			} else {
				relFiles[i] = rel
			}
		}
		result = append(result, relFiles)
	}
	return result, nil
}

func isGoFile(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasPrefix(name, ".")
}

func isExcludedDir(dir string, excludeDirs []string) bool {
	for _, ex := range excludeDirs {
		if dir == ex || strings.HasPrefix(dir, ex+string(filepath.Separator)) {
			return true
		}
	}
	return false
}
