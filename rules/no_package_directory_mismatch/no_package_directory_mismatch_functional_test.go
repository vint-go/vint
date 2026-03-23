package no_package_directory_mismatch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_package_directory_mismatch"
)

var emptyIgnoreConfig = &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignore-directories": []any{}}}}

func TestNoPackageDirectoryMismatch(t *testing.T) {
	functional_test_helpers.TestRule(t, "good/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "bad/bad", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "maincmd/main", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "mixed/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "mixed/bad", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "go-good/good1", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "go-good/good2", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "go-good/bad", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test normalization cases
	functional_test_helpers.TestRule(t, "normalization/fo-ob_ar/foobar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/foo_bar/foobar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/foo.b_ar/foobar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-foo_bar/foo_bar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-foo_bar/foobar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-foo-bar/foo_bar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-foo-bar/foobar", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Edge cases that are not logic, but we decided to ignore to have a simpler implementation
	functional_test_helpers.TestRule(t, "normalization/go-od/go_od", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-od/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "normalization/go-od/od", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test version directories (v1, v2, etc.)
	functional_test_helpers.TestRule(t, "api/v1/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "api/V1/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "api/V1/api_test", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "api/v1v/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "api/v2/v2", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test internal directory variations
	functional_test_helpers.TestRule(t, "internal/good/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "internal/bad/bad", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "internal/any", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "internal/api/v1/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "internal/api/v2/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "internal/v1/api", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test handling of testfiles
	functional_test_helpers.TestRule(t, "test/good_test", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "test/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "test/bad_test", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "test/bad", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "test/main_test", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test handling of root directories with go.mod
	functional_test_helpers.TestRule(t, "rootdir/good/client", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
	functional_test_helpers.TestRule(t, "rootdir/bad/client", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)

	// Test handling of root directories with .git
	mkdirTempDotGit(t, "rootdir/withgit")
	functional_test_helpers.TestRule(t, "rootdir/withgit/client", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, emptyIgnoreConfig)
}

func TestNoPackageDirectoryMismatchWithDefaultConfig(t *testing.T) {
	// Test with default configuration (should ignore testdata directories by default)
	functional_test_helpers.TestRule(t, "testdata_dir/ignored", &no_package_directory_mismatch.PackageDirectoryMismatchRule{})
}

func TestNoPackageDirectoryMismatchWithTestDirectories(t *testing.T) {
	// Test with new format that excludes specific test directories
	config := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignore-directories": []any{"testinfo", "testutils"}}}}

	functional_test_helpers.TestRule(t, "testinfo/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, config)
	functional_test_helpers.TestRule(t, "testutils/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, config)
}

func TestNoPackageDirectoryMismatchWithMultipleDirectories(t *testing.T) {
	// Test with new named argument format that excludes both "testutils" and "testinfo"
	config := &lint.RuleConfig{Arguments: lint.Arguments{map[string]any{"ignoreDirectories": []any{"testutils", "testinfo"}}}}

	functional_test_helpers.TestRule(t, "testutils/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, config)
	functional_test_helpers.TestRule(t, "testinfo/good", &no_package_directory_mismatch.PackageDirectoryMismatchRule{}, config)
}

// mkdirTempDotGit adds a temporary .git directory to the given root directory in testdata.
// We can't commit .git directly because of the Git restrictions.
func mkdirTempDotGit(t *testing.T, root string) {
	t.Helper()

	baseDir := filepath.Join("testdata", root)
	dir, err := filepath.Abs(baseDir)
	if err != nil {
		t.Fatalf("Failed to resolve abs path: %v", err)
	}

	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o750); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	t.Cleanup(func() {
		err = os.RemoveAll(gitDir)
		if err != nil {
			t.Fatalf("Failed to remove .git directory: %v", err)
		}
	})
}
