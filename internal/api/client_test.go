package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Authorization = %q, want 'Bearer test-token'", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/v1/accounts" {
			t.Errorf("Path = %q, want '/v1/accounts'", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"accounts": []any{}})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	resp, err := client.Get(context.Background(), "/v1/accounts")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q, want 'application/json'", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	resp, err := client.Post(context.Background(), "/v1/chats/123/messages", map[string]string{"text": "hello"})
	if err != nil {
		t.Fatalf("Post() error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}
}
