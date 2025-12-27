package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDebugLogging(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create client with debug enabled
	client := NewClient(server.URL, "test-token", WithDebug(true))

	// Make a request
	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	_ = resp.Body.Close()

	// Restore stderr and read output
	_ = w.Close()
	os.Stderr = oldStderr

	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Verify debug output
	if !strings.Contains(output, "→ GET") {
		t.Errorf("Expected request debug output, got: %s", output)
	}
	if !strings.Contains(output, "← 200") {
		t.Errorf("Expected response debug output, got: %s", output)
	}
	if strings.Contains(output, "test-token") {
		t.Errorf("Debug output should not contain the auth token")
	}
}

func TestDebugDisabled(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Create client without debug (default)
	client := NewClient(server.URL, "test-token")

	// Make a request
	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	_ = resp.Body.Close()

	// Restore stderr and read output
	_ = w.Close()
	os.Stderr = oldStderr

	outputBytes, _ := io.ReadAll(r)
	output := string(outputBytes)

	// Verify NO debug output
	if strings.Contains(output, "→ GET") {
		t.Errorf("Expected no debug output, got: %s", output)
	}
}
