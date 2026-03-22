package fixtures

import "os"

// Invalid: 0666 exceeds custom max 0644
func openFileTooPermissive() (*os.File, error) {
	return os.OpenFile("data.txt", os.O_CREATE|os.O_WRONLY, 0666) // MATCH /file permission 0666 is more permissive than 0644/
}

// Invalid: 0755 exceeds custom max 0644
func chmodTooPermissive() error {
	return os.Chmod("script.sh", 0755) // MATCH /file permission 0755 is more permissive than 0644/
}

// Valid: exactly at threshold
func openFileAtThreshold() (*os.File, error) {
	return os.OpenFile("readme.txt", os.O_CREATE|os.O_WRONLY, 0644)
}

// Valid: below threshold
func openFileRestricted() (*os.File, error) {
	return os.OpenFile("secret.key", os.O_CREATE|os.O_WRONLY, 0600)
}

// Valid: 0600 is well below custom max 0644 (would fail with default 0600 threshold)
func chmodDefault() error {
	return os.Chmod("config.yml", 0600)
}
