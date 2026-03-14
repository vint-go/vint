package fixtures

import (
	"io/ioutil"
	"os"
)

func saveConfigWorldReadable(data []byte) error {
	return os.WriteFile("config.yml", data, 0644) // MATCH /file permission 0644 is more permissive than 0600/
}

func saveKeyOverlyPermissive(key []byte) error {
	return os.WriteFile("private.key", key, 0666) // MATCH /file permission 0666 is more permissive than 0600/
}

func saveWorldReadWrite(data []byte) error {
	return os.WriteFile("data.txt", data, 0777) // MATCH /file permission 0777 is more permissive than 0600/
}

func saveGroupReadable(data []byte) error {
	return os.WriteFile("shared.conf", data, 0640) // MATCH /file permission 0640 is more permissive than 0600/
}

func saveWithIoutilWorldReadable(data []byte) error {
	return ioutil.WriteFile("legacy.conf", data, 0644) // MATCH /file permission 0644 is more permissive than 0600/
}

func saveConfigOwnerOnly(data []byte) error {
	return os.WriteFile("config.yml", data, 0600)
}

func saveKeyReadOnly(key []byte) error {
	return os.WriteFile("private.key", key, 0400)
}

func saveWithIoutilOwnerOnly(data []byte) error {
	return ioutil.WriteFile("legacy.conf", data, 0600)
}

func saveNoPermission(data []byte) error {
	return os.WriteFile("empty.dat", data, 0000)
}
