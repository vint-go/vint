package fixtures

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
)

// Invalid: Direct taint from HTTP request to log.Printf
func handlerLogPrintf(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	// Taint flows from request to log output
	log.Printf("Login attempt for user: %s", username) // MATCH /potential log injection: tainted data from user input flows to log.Printf/
}

// Invalid: Direct taint from HTTP request query to slog.Info with raw string arg
func handlerSlogInfo(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("search")
	// Unsanitized input in log message
	slog.Info("Search performed", "query", input) // MATCH /potential log injection: tainted data from user input flows to slog.Info/
}

// Invalid: Taint flows through fmt.Sprintf to log.Println
func handlerSprintfLog(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	msg := fmt.Sprintf("User %s logged in", name)
	log.Println(msg) // MATCH /potential log injection: tainted data from user input flows to log.Println/
}

// Invalid: Taint flows through string concatenation to log.Print
func handlerConcatLog(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	logMsg := "Action: " + action
	log.Print(logMsg) // MATCH /potential log injection: tainted data from user input flows to log.Print/
}

// Invalid: Interprocedural - taint flows through helper function
func logUserAction(action string) {
	log.Printf("User action: %s", action)
}

func handlerInterprocedural(w http.ResponseWriter, r *http.Request) {
	userAction := r.FormValue("action")
	logUserAction(userAction) // MATCH /potential log injection: tainted data from user input flows to log.Printf/
}

// Valid: Sanitize control characters before logging
func handlerSanitized(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	// Sanitize control characters before logging
	clean := strings.NewReplacer("\n", "", "\r", "", "\x00", "").Replace(username)
	log.Printf("Login attempt for user: %s", clean)
}

// Valid: Structured logging with slog typed attributes
func handlerSlogTyped(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("search")
	// Structured logging with proper encoding
	slog.Info("Search performed",
		slog.String("query", input),
		slog.String("remote_addr", r.RemoteAddr),
	)
}

// Valid: No user input in log message
func handlerNoUserInput() {
	log.Printf("Request received at %s", "localhost:8080")
}

// Valid: Logging a hardcoded constant
func safeLogging() {
	log.Println("Application started successfully")
}
