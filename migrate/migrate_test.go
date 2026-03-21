package migrate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/strowk/vint/migrate"

	// Register linter migrators.
	_ "github.com/strowk/vint/migrate/migrator/linters/bodyclose"
	_ "github.com/strowk/vint/migrate/migrator/linters/copyloopvar"
	_ "github.com/strowk/vint/migrate/migrator/linters/depguard"
	_ "github.com/strowk/vint/migrate/migrator/linters/dogsled"
	_ "github.com/strowk/vint/migrate/migrator/linters/dupl"
	_ "github.com/strowk/vint/migrate/migrator/linters/err113"
	_ "github.com/strowk/vint/migrate/migrator/linters/errcheck"
	_ "github.com/strowk/vint/migrate/migrator/linters/errorlint"
	_ "github.com/strowk/vint/migrate/migrator/linters/funlen"
	_ "github.com/strowk/vint/migrate/migrator/linters/gocheckcompilerdirectives"
	_ "github.com/strowk/vint/migrate/migrator/linters/gochecknoinits"
	_ "github.com/strowk/vint/migrate/migrator/linters/goconst"
	_ "github.com/strowk/vint/migrate/migrator/linters/gocritic"
	_ "github.com/strowk/vint/migrate/migrator/linters/gocyclo"
	_ "github.com/strowk/vint/migrate/migrator/linters/goprintffuncname"
	_ "github.com/strowk/vint/migrate/migrator/linters/gosec"
	_ "github.com/strowk/vint/migrate/migrator/linters/govet"
	_ "github.com/strowk/vint/migrate/migrator/linters/ineffassign"
	_ "github.com/strowk/vint/migrate/migrator/linters/intrange"
	_ "github.com/strowk/vint/migrate/migrator/linters/lll"
	_ "github.com/strowk/vint/migrate/migrator/linters/misspell"
	_ "github.com/strowk/vint/migrate/migrator/linters/mnd"
	_ "github.com/strowk/vint/migrate/migrator/linters/nakedret"
	_ "github.com/strowk/vint/migrate/migrator/linters/noctx"
	_ "github.com/strowk/vint/migrate/migrator/linters/nolintlint"
	_ "github.com/strowk/vint/migrate/migrator/linters/sloglint"
	_ "github.com/strowk/vint/migrate/migrator/linters/staticcheck"
	_ "github.com/strowk/vint/migrate/migrator/linters/unconvert"
	_ "github.com/strowk/vint/migrate/migrator/linters/unparam"
	_ "github.com/strowk/vint/migrate/migrator/linters/unused"
	_ "github.com/strowk/vint/migrate/migrator/linters/whitespace"
)

func TestMigrateConfigExamples(t *testing.T) {
	examples := findExamples(t, "examples")

	for _, ex := range examples {
		t.Run(ex.name, func(t *testing.T) {
			configPath := filepath.Join(ex.dir, ".golangci.yml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				t.Skipf("no .golangci.yml in %s", ex.dir)
			}

			result, err := migrate.Migrate(configPath, ex.dir)
			if err != nil {
				t.Fatalf("Migrate() error: %v", err)
			}

			// Check vint.yaml output.
			expectedYAMLPath := filepath.Join(ex.dir, "expected.vint.yaml")
			if _, err := os.Stat(expectedYAMLPath); err == nil {
				expected, err := os.ReadFile(expectedYAMLPath)
				if err != nil {
					t.Fatalf("read expected.vint.yaml: %v", err)
				}

				if normalizeYAML(result.VintYAML) != normalizeYAML(string(expected)) {
					t.Errorf("vint.yaml mismatch:\n--- expected ---\n%s\n--- got ---\n%s",
						string(expected), result.VintYAML)
				}
			}

			// Check converted Go files.
			for filePath, gotContent := range result.ConvertedFiles {
				baseName := filepath.Base(filePath)
				expectedPath := filepath.Join(ex.dir, baseName+".expected")
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Logf("no expected file for %s, skipping comparison", baseName)
					t.Logf("got:\n%s", gotContent)
					continue
				}

				expected, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("read %s: %v", expectedPath, err)
				}

				if normalizeSource(gotContent) != normalizeSource(string(expected)) {
					t.Errorf("%s mismatch:\n--- expected ---\n%s\n--- got ---\n%s",
						baseName, string(expected), gotContent)
				}
			}
		})
	}
}

type example struct {
	name string
	dir  string
}

func findExamples(t *testing.T, root string) []example {
	t.Helper()
	var examples []example

	linters, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("read examples dir: %v", err)
	}

	for _, linter := range linters {
		if !linter.IsDir() {
			continue
		}
		cases, err := os.ReadDir(filepath.Join(root, linter.Name()))
		if err != nil {
			t.Fatalf("read linter dir %s: %v", linter.Name(), err)
		}
		for _, c := range cases {
			if !c.IsDir() {
				continue
			}
			examples = append(examples, example{
				name: linter.Name() + "/" + c.Name(),
				dir:  filepath.Join(root, linter.Name(), c.Name()),
			})
		}
	}

	return examples
}

func normalizeYAML(s string) string {
	return strings.TrimSpace(s)
}

func normalizeSource(s string) string {
	return strings.TrimSpace(s)
}
