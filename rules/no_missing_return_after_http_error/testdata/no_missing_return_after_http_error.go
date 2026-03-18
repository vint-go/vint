package fixtures

import "net/http"

// Invalid: http.Error without a return.
func badHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest) // MATCH /http.Error call is not followed by a return statement/
	}
	w.Write([]byte("hello"))
}

// Invalid: http.Error as the last statement in a block (no return follows).
func badHandlerLastStmt(w http.ResponseWriter, r *http.Request) {
	var err error
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError) // MATCH /http.Error call is not followed by a return statement/
	}
}

// Invalid: http.Error followed by another non-return statement.
func badHandlerFollowedByLog(w http.ResponseWriter, r *http.Request) {
	var err error
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden) // MATCH /http.Error call is not followed by a return statement/
		_ = err
	}
}

// Valid: http.Error followed by a return.
func goodHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	w.Write([]byte("hello"))
}

// Valid: http.Error at the end of a function that is the last statement and no other writes.
func goodHandlerTerminal(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
	return
}
