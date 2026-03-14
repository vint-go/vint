package fixtures

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// Invalid: Direct taint from HTTP request to exec.Command via function call
func handler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// Taint flows from HTTP request to command execution
	output := runCommand(filename)  // MATCH /potential command injection: tainted data from user input flows to exec.Command/
	w.Write(output)
}

func runCommand(file string) []byte {
	out, _ := exec.Command("cat", file).Output()
	return out
}

// Invalid: User input passed to shell (detected at call site)
func processInput(input string) error {
	cmd := exec.Command("sh", "-c", input)
	return cmd.Run()
}

func handlerProcessInput(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("cmd")
	processInput(userInput) // MATCH /potential command injection: tainted data from user input flows to exec.Command/
}

// Invalid: Taint flows through string concatenation
func handlerConcat(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	cmd := "echo " + name
	exec.Command("sh", "-c", cmd) // MATCH /potential command injection: tainted data from user input flows to exec.Command/
}

// Invalid: Taint flows through fmt.Sprintf
func handlerSprintf(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("host")
	arg := fmt.Sprintf("ping %s", host)
	exec.Command("sh", "-c", arg) // MATCH /potential command injection: tainted data from user input flows to exec.Command/
}

// Valid: Hardcoded constant command with no user input
func safeCommand() error {
	return exec.Command("/usr/bin/process", "safe-arg").Run()
}

// Valid: Validate against allowlist before using in command
func handlerSafe(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	// Validate against allowlist
	if !isAllowedFile(filename) {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}
	// Read file directly instead of using command
	data, err := os.ReadFile(filepath.Join("/safe/dir", filepath.Base(filename)))
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func isAllowedFile(name string) bool {
	return name == "allowed.txt"
}

// Valid: Use constant command with validated argument
func processFile(filename string) error {
	clean := filepath.Base(filename)
	if !regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`).MatchString(clean) {
		return fmt.Errorf("invalid filename")
	}
	return exec.Command("/usr/bin/process", clean).Run()
}
