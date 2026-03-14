package no_unsafe_package_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unsafe_package"
)

func TestNoUnsafePackage(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unsafe_package", &no_unsafe_package.NoUnsafePackageRule{})
}
