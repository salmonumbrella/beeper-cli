package api

import (
	"sync"
	"time"
)

// circuitBreaker implements the circuit breaker pattern for resilience.
// It opens after CircuitBreakerThreshold consecutive failures and
// auto-resets after CircuitBreakerResetTime to allow retry.
type circuitBreaker struct {
	mu          sync.Mutex
	failures    int
	lastFailure time.Time
	open        bool
}

// recordSuccess resets the failure count and closes the circuit.
func (cb *circuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.open = false
}

// recordFailure increments the failure count and opens the circuit
// if the threshold is reached. Returns true if the circuit just opened.
func (cb *circuitBreaker) recordFailure() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFailure = time.Now()
	if cb.failures >= CircuitBreakerThreshold {
		cb.open = true
		return true
	}
	return false
}

// isOpen returns true if the circuit is open (requests should fail fast).
// If the reset time has passed, it automatically resets the circuit.
func (cb *circuitBreaker) isOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if !cb.open {
		return false
	}
	// Check if reset time has passed
	if time.Since(cb.lastFailure) > CircuitBreakerResetTime {
		cb.open = false
		cb.failures = 0
		return false
	}
	return true
}
