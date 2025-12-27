// internal/testutil/testutil_test.go
package testutil

import (
	"io"
	"net/http"
	"testing"
)

func TestNewMockServer(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}

	server := NewMockServer(t, handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestNewMockServerWithRoutes(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/api/v1/users": func(w http.ResponseWriter, r *http.Request) {
			JSONResponse(w, http.StatusOK, `{"users":[]}`)
		},
		"/api/v1/status": func(w http.ResponseWriter, r *http.Request) {
			JSONResponse(w, http.StatusOK, `{"status":"healthy"}`)
		},
	}

	server := NewMockServerWithRoutes(t, routes)
	defer server.Close()

	// Test /api/v1/users route
	resp, err := http.Get(server.URL + "/api/v1/users")
	if err != nil {
		t.Fatalf("request to /api/v1/users failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"users":[]}` {
		t.Errorf("unexpected body: %s", body)
	}

	// Test /api/v1/status route
	resp2, err := http.Get(server.URL + "/api/v1/status")
	if err != nil {
		t.Fatalf("request to /api/v1/status failed: %v", err)
	}
	defer func() { _ = resp2.Body.Close() }()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp2.StatusCode)
	}

	body2, _ := io.ReadAll(resp2.Body)
	if string(body2) != `{"status":"healthy"}` {
		t.Errorf("unexpected body: %s", body2)
	}
}

func TestJSONResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, http.StatusCreated, `{"id":123}`)
	}

	server := NewMockServer(t, handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"id":123}` {
		t.Errorf("unexpected body: %s", body)
	}
}
