package fixtures

import "os"

// Invalid: 0666 exceeds custom max 0644
func writeFileTooPermissive() error {
	return os.WriteFile("data.txt", []byte("hello"), 0666) // MATCH /file permission 0666 is more permissive than 0644/
}

// Invalid: 0777 exceeds custom max 0644
func writeFileWayTooPermissive() error {
	return os.WriteFile("script.sh", []byte("#!/bin/bash"), 0777) // MATCH /file permission 0777 is more permissive than 0644/
}

// Valid: exactly at threshold
func writeFileAtThreshold() error {
	return os.WriteFile("readme.txt", []byte("readme"), 0644)
}

// Valid: below threshold
func writeFileRestricted() error {
	return os.WriteFile("secret.key", []byte("key"), 0600)
}
