// internal/api/client_circuit_test.go
package api

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/salmonumbrella/beeper-cli/internal/testutil"
)

func TestClientCircuitBreakerOpens(t *testing.T) {
	var attempts atomic.Int32

	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	})

	client := NewClient(server.URL, "test-token")

	// Make enough requests to open the circuit
	// Each GET request retries once (Max5xxRetries=1), so failures accumulate
	for i := 0; i < CircuitBreakerThreshold; i++ {
		_, _ = client.Get(context.Background(), "/test")
	}

	// Next request should fail immediately with circuit open error
	_, err := client.Get(context.Background(), "/test")
	if err == nil {
		t.Error("expected circuit breaker error")
	}
	if err.Error() != "circuit breaker open: API experiencing issues, retry later" {
		t.Errorf("unexpected error: %v", err)
	}
}
