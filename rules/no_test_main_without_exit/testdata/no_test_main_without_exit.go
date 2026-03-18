package fixtures

import (
	"os"
	"testing"
)

// Invalid: TestMain without os.Exit - test failures will be hidden.
func TestMain(m *testing.M) { // MATCH /TestMain should call os.Exit to set exit code, otherwise test failures will be hidden/
	// setup
	m.Run()
	// teardown
}

// Valid: not named TestMain, so no check applies.
func TestSomething(m *testing.M) {
	m.Run()
}

// Valid: function with different signature (not *testing.M param).
func helperFunc() {
	// not a TestMain
}

// helper to use os import
var _ = os.Exit
