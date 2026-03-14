package fixtures

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Invalid: Direct taint from HTTP request to os.ReadFile via function call
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// Taint flows from request to file read
	data := readUserFile(filename) // MATCH /potential path traversal: tainted data from user input flows to os.ReadFile/
	w.Write(data)
}

func readUserFile(name string) []byte {
	path := filepath.Join("/uploads", name)
	data, _ := os.ReadFile(path)
	return data
}

// Invalid: Taint flows through filepath.Join directly
func handlerDirect(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	path := filepath.Join("/uploads", filename)
	os.ReadFile(path) // MATCH /potential path traversal: tainted data from user input flows to os.ReadFile/
}

// Invalid: Taint flows through string concatenation to os.Open
func handlerConcat(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := "/data/" + name
	os.Open(path) // MATCH /potential path traversal: tainted data from user input flows to os.Open/
}

// Invalid: Taint flows through fmt.Sprintf to os.Create
func handlerSprintf(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	path := fmt.Sprintf("/uploads/%s", name)
	os.Create(path) // MATCH /potential path traversal: tainted data from user input flows to os.Create/
}

// Valid: Hardcoded constant path with no user input
func safeReadFile() error {
	_, err := os.ReadFile("/etc/config.txt")
	return err
}

// Valid: Validate path stays within base directory
func downloadHandlerSafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	data, err := readUserFileSafe(filename)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Write(data)
}

func readUserFileSafe(name string) ([]byte, error) {
	baseDir := "/uploads"
	absPath, err := filepath.Abs(filepath.Join(baseDir, name))
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(absPath, baseDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("path traversal detected")
	}
	return os.ReadFile(absPath)
}

// Valid: Static file operations without user input
func writeConfig() error {
	return os.WriteFile("/etc/app/config.json", []byte("{}"), 0644)
}
