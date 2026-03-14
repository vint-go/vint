package fixtures

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Invalid: file path from variable parameter - potential directory traversal
func readFileVariable(path string) ([]byte, error) {
	return os.ReadFile(path) // MATCH /potential file inclusion via variable: path passed to os.ReadFile is not a hardcoded constant/
}

// Invalid: user-controlled path in os.Open
func openFileVariable(userPath string) (*os.File, error) {
	return os.Open(userPath) // MATCH /potential file inclusion via variable: path passed to os.Open is not a hardcoded constant/
}

// Invalid: concatenation with user input
func readUserFile(filename string) ([]byte, error) {
	path := "/data/uploads/" + filename
	return os.ReadFile(path) // MATCH /potential file inclusion via variable: path passed to os.ReadFile is not a hardcoded constant/
}

// Invalid: os.OpenFile with variable path
func openFileVariableFlags(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_RDONLY, 0644) // MATCH /potential file inclusion via variable: path passed to os.OpenFile is not a hardcoded constant/
}

// Invalid: os.Create with variable path
func createFileVariable(name string) (*os.File, error) {
	return os.Create(name) // MATCH /potential file inclusion via variable: path passed to os.Create is not a hardcoded constant/
}

// Invalid: ioutil.ReadFile with variable path
func readFileIoutil(path string) ([]byte, error) {
	return ioutil.ReadFile(path) // MATCH /potential file inclusion via variable: path passed to ioutil.ReadFile is not a hardcoded constant/
}

// Valid: hardcoded constant path
func readConfig() ([]byte, error) {
	return os.ReadFile("/etc/myapp/config.yaml")
}

// Valid: hardcoded constant in os.Open
func openConfig() (*os.File, error) {
	return os.Open("/etc/myapp/config.yaml")
}

// Valid: path sanitized via filepath.Clean
func readFileSanitized(basePath, filename string) ([]byte, error) {
	cleanPath := filepath.Clean(filepath.Join(basePath, filename))
	return os.ReadFile(cleanPath)
}

// Valid: path sanitized via filepath.Rel
func readFileRel(basePath, targetPath string) ([]byte, error) {
	relPath, _ := filepath.Rel(basePath, targetPath)
	return os.ReadFile(relPath)
}

// Valid: path sanitized via filepath.EvalSymlinks
func readFileEvalSymlinks(path string) ([]byte, error) {
	evalPath, _ := filepath.EvalSymlinks(path)
	return os.ReadFile(evalPath)
}

// Valid: hardcoded constant in os.Create
func createConfigFile() (*os.File, error) {
	return os.Create("/tmp/myapp/output.txt")
}

// Valid: hardcoded constant in os.OpenFile
func openConfigFile() (*os.File, error) {
	return os.OpenFile("/etc/myapp/data.txt", os.O_RDONLY, 0644)
}
