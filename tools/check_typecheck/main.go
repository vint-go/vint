// check_typecheck scans rule source files for calls to .Pkg.TypeCheck()
// and reports any that don't also declare RequiresTypecheck() bool.
//
// Usage: go run ./tools/check_typecheck [dirs...]
//
// If no directories are given, it defaults to ./rule and ./rules.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dirs := os.Args[1:]
	if len(dirs) == 0 {
		dirs = []string{"rule", "rules"}
	}

	var bad []string
	for _, dir := range dirs {
		found, err := check(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error scanning %s: %v\n", dir, err)
			os.Exit(1)
		}
		bad = append(bad, found...)
	}

	if len(bad) == 0 {
		fmt.Println("OK: all rules that call Pkg.TypeCheck() also declare RequiresTypecheck().")
		return
	}

	fmt.Fprintf(os.Stderr, "FAIL: %d rule file(s) call Pkg.TypeCheck() without declaring RequiresTypecheck():\n", len(bad))
	for _, f := range bad {
		fmt.Fprintf(os.Stderr, "  %s\n", f)
	}
	os.Exit(1)
}

func check(dir string) ([]string, error) {
	var bad []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			return err
		}
		src := string(data)

		callsTypeCheck := strings.Contains(src, ".Pkg.TypeCheck()")
		declaresInterface := strings.Contains(src, "RequiresTypecheck()")

		if callsTypeCheck && !declaresInterface {
			bad = append(bad, path)
		}
		return nil
	})
	return bad, err
}
