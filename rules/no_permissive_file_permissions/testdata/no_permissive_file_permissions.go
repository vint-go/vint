package fixtures

import "os"

func createFileWorldReadable() (*os.File, error) {
	return os.OpenFile("config.yml", os.O_CREATE|os.O_WRONLY, 0644) // MATCH /file permission 0644 is more permissive than 0600/
}

func chmodOverlyPermissive() error {
	return os.Chmod("secret.key", 0666) // MATCH /file permission 0666 is more permissive than 0600/
}

func createExecutableWorldExecutable() (*os.File, error) {
	return os.OpenFile("script.sh", os.O_CREATE|os.O_WRONLY, 0755) // MATCH /file permission 0755 is more permissive than 0600/
}

func chmodWorldReadWrite() error {
	return os.Chmod("data.txt", 0777) // MATCH /file permission 0777 is more permissive than 0600/
}

func createFileOwnerOnly() (*os.File, error) {
	return os.OpenFile("config.yml", os.O_CREATE|os.O_WRONLY, 0600)
}

func chmodOwnerOnly() error {
	return os.Chmod("secret.key", 0600)
}

func createFileReadOnly() (*os.File, error) {
	return os.OpenFile("readme.txt", os.O_RDONLY, 0400)
}

func chmodMinimal() error {
	return os.Chmod("private.key", 0400)
}
