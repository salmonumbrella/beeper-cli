// internal/testutil/testutil.go
package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewMockServer creates a test HTTP server with the given handler.
// The server is automatically closed when the test completes.
func NewMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// NewMockServerWithRoutes creates a test server with multiple route handlers.
func NewMockServerWithRoutes(t *testing.T, routes map[string]http.HandlerFunc) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	for path, handler := range routes {
		mux.HandleFunc(path, handler)
	}
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server
}

// JSONResponse is a helper to write JSON responses in tests.
func JSONResponse(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
