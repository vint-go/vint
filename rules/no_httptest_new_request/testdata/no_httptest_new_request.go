package fixtures

import (
	"context"
	"net/http"
	"net/http/httptest"
)

func handler(w http.ResponseWriter, r *http.Request) {}

func badHttptestNewRequest() {
	req := httptest.NewRequest(http.MethodGet, "/path", nil) // MATCH /httptest.NewRequest does not accept a context; use httptest.NewRequestWithContext instead/
	w := httptest.NewRecorder()
	handler(w, req)
}

func goodHttptestNewRequestWithContext() {
	ctx := context.Background()
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/path", nil)
	w := httptest.NewRecorder()
	handler(w, req)
}
