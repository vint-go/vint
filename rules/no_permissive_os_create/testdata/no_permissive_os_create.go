package fixtures

import "os"

func createFile() (*os.File, error) {
	// os.Create uses default 0666 permissions
	return os.Create("sensitive-data.txt") // MATCH /os.Create uses default 0666 permissions; use os.OpenFile with explicit permissions instead/
}

func writeSecret(secret string) error {
	f, err := os.Create("secret.key") // MATCH /os.Create uses default 0666 permissions; use os.OpenFile with explicit permissions instead/
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(secret)
	return err
}

func createFileExplicit() (*os.File, error) {
	// Explicit restrictive permissions
	return os.OpenFile("sensitive-data.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

func writeSecretSafe(secret string) error {
	return os.WriteFile("secret.key", []byte(secret), 0600)
}
