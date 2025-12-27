// internal/api/client_retry_test.go
package api

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/salmonumbrella/beeper-cli/internal/testutil"
)

func TestClientRetryOn429(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		count := attempts.Add(1)
		if count < 3 {
			w.Header().Set("Retry-After", "0") // Immediate retry for test
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		testutil.JSONResponse(w, http.StatusOK, `{"status":"ok"}`)
	})

	client := NewClient(server.URL, "test-token")
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestClientRetryOn5xx(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		count := attempts.Add(1)
		if count <= 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		testutil.JSONResponse(w, http.StatusOK, `{"status":"ok"}`)
	})

	client := NewClient(server.URL, "test-token")
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if attempts.Load() != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestClientNo5xxRetryOnPost(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(server.URL, "test-token")
	resp, err := client.Post(context.Background(), "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// POST is not idempotent, so no retry on 5xx
	if attempts.Load() != 1 {
		t.Errorf("expected 1 attempt (no retry for POST), got %d", attempts.Load())
	}
}

func TestClientRetryRespectsContextCancellation(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Retry-After", "10")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	client := NewClient(server.URL, "test-token")
	_, err := client.Get(ctx, "/test")

	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got %v", err)
	}
	if attempts.Load() != 1 {
		t.Errorf("expected 1 attempt before context cancelled, got %d", attempts.Load())
	}
}

func TestClientRetryExhaustsRetries(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	})

	client := NewClient(server.URL, "test-token")
	resp, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// After MaxRateLimitRetries (3), returns the 429 response
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", resp.StatusCode)
	}
	if attempts.Load() != int32(MaxRateLimitRetries) {
		t.Errorf("expected %d attempts, got %d", MaxRateLimitRetries, attempts.Load())
	}
}
