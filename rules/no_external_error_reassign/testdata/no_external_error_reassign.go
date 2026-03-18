package fixtures

import (
	"errors"
	"io"
	"net/http"
)

func invalidReassignments() {
	io.EOF = nil                                      // MATCH /reassigning external error variable io.EOF/
	http.ErrServerClosed = errors.New("custom error") // MATCH /reassigning external error variable http.ErrServerClosed/
}

func validUsages() {
	if errors.Is(nil, io.EOF) {
		// compare, don't reassign
	}

	err := io.EOF
	_ = err

	if err == http.ErrServerClosed {
		// compare, don't reassign
	}
}
