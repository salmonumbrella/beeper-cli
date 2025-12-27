package api

import (
	"testing"
	"time"
)

func TestCircuitBreakerOpensAfterThreshold(t *testing.T) {
	cb := &circuitBreaker{}

	// Record failures up to threshold
	for i := 0; i < CircuitBreakerThreshold-1; i++ {
		opened := cb.recordFailure()
		if opened {
			t.Errorf("circuit opened too early at failure %d", i+1)
		}
		if cb.isOpen() {
			t.Errorf("circuit should not be open yet")
		}
	}

	// This failure should open the circuit
	opened := cb.recordFailure()
	if !opened {
		t.Error("circuit should have opened")
	}
	if !cb.isOpen() {
		t.Error("circuit should be open")
	}
}

func TestCircuitBreakerResetsOnSuccess(t *testing.T) {
	cb := &circuitBreaker{}

	// Record some failures
	for i := 0; i < CircuitBreakerThreshold-1; i++ {
		cb.recordFailure()
	}

	// Success should reset
	cb.recordSuccess()

	if cb.failures != 0 {
		t.Errorf("failures should be 0, got %d", cb.failures)
	}
}

func TestCircuitBreakerResetsAfterTimeout(t *testing.T) {
	cb := &circuitBreaker{}

	// Open the circuit
	for i := 0; i < CircuitBreakerThreshold; i++ {
		cb.recordFailure()
	}

	if !cb.isOpen() {
		t.Fatal("circuit should be open")
	}

	// Simulate time passing by setting lastFailure in the past
	cb.mu.Lock()
	cb.lastFailure = time.Now().Add(-CircuitBreakerResetTime - time.Second)
	cb.mu.Unlock()

	// Circuit should now be closed (reset)
	if cb.isOpen() {
		t.Error("circuit should be closed after reset time")
	}

	// Verify internal state was reset
	cb.mu.Lock()
	if cb.failures != 0 {
		t.Errorf("failures should be 0 after reset, got %d", cb.failures)
	}
	if cb.open {
		t.Error("open flag should be false after reset")
	}
	cb.mu.Unlock()
}

func TestCircuitBreakerStaysOpenBeforeTimeout(t *testing.T) {
	cb := &circuitBreaker{}

	// Open the circuit
	for i := 0; i < CircuitBreakerThreshold; i++ {
		cb.recordFailure()
	}

	// Circuit should still be open (not enough time passed)
	if !cb.isOpen() {
		t.Error("circuit should still be open before reset time")
	}
}
