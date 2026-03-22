package fixtures

import "os"

// Invalid: 0777 exceeds custom max 0700
func mkdirTooPermissive() error {
	return os.Mkdir("/tmp/myapp", 0777) // MATCH /directory permission 0777 is more permissive than 0700/
}

// Invalid: 0750 exceeds custom max 0700
func mkdirAllGroupAccess() error {
	return os.MkdirAll("/var/myapp/data", 0750) // MATCH /directory permission 0750 is more permissive than 0700/
}

// Valid: exactly at threshold
func mkdirExactThreshold() error {
	return os.Mkdir("/tmp/safe", 0700)
}

// Valid: below threshold
func mkdirRestricted() error {
	return os.MkdirAll("/tmp/private", 0500)
}
