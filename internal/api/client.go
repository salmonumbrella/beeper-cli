package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	DefaultBaseURL     = "http://localhost:23373"
	DefaultHTTPTimeout = 30 * time.Second

	// MaxRateLimitRetries is the maximum number of retries on 429 responses.
	MaxRateLimitRetries = 3

	// RateLimitBaseDelay is the initial delay for rate limit exponential backoff.
	RateLimitBaseDelay = 1 * time.Second

	// Max5xxRetries is the maximum retries for server errors on idempotent requests.
	Max5xxRetries = 1

	// ServerErrorRetryDelay is the delay before retrying on 5xx errors.
	ServerErrorRetryDelay = 1 * time.Second

	// CircuitBreakerThreshold is consecutive 5xx errors to open circuit.
	CircuitBreakerThreshold = 5

	// CircuitBreakerResetTime is how long before trying again.
	CircuitBreakerResetTime = 30 * time.Second
)

type Client struct {
	baseURL        string
	token          string
	httpClient     *http.Client
	debug          bool
	circuitBreaker *circuitBreaker
}

type ClientOption func(*Client)

func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

func NewClient(baseURL, token string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: DefaultHTTPTimeout,
		},
		circuitBreaker: &circuitBreaker{},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	return c.doWithRetry(ctx, req)
}

// doWithRetry handles rate limiting (429) and server errors (5xx) with retries.
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Check circuit breaker at start - fail fast if open
	if c.circuitBreaker.isOpen() {
		return nil, fmt.Errorf("circuit breaker open: API experiencing issues, retry later")
	}

	var rateLimitRetries int
	var serverErrorRetries int

	for {
		if c.debug {
			fmt.Fprintf(os.Stderr, "→ %s %s\n", req.Method, req.URL)
			for k, v := range req.Header {
				if k != "Authorization" { // Don't log token
					fmt.Fprintf(os.Stderr, "  %s: %s\n", k, v[0])
				}
			}
		}

		resp, err := c.httpClient.Do(req)

		if c.debug && resp != nil {
			fmt.Fprintf(os.Stderr, "← %d %s\n", resp.StatusCode, resp.Status)
		}

		if err != nil {
			return nil, err
		}

		// Handle 429 Too Many Requests with exponential backoff
		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitRetries++
			if rateLimitRetries >= MaxRateLimitRetries {
				return resp, nil
			}

			// Close the response body before retrying
			_ = resp.Body.Close()

			delay := c.parseRetryAfter(resp.Header.Get("Retry-After"), rateLimitRetries)

			if c.debug {
				fmt.Fprintf(os.Stderr, "  rate limited, retry %d/%d in %v\n", rateLimitRetries, MaxRateLimitRetries, delay)
			}

			// Replay request body if needed
			if err := c.replayRequestBody(req); err != nil {
				return nil, fmt.Errorf("failed to replay request body: %w", err)
			}

			if err := c.sleep(ctx, delay); err != nil {
				return nil, err
			}
			continue
		}

		// Handle 5xx server errors with single retry for idempotent methods
		if resp.StatusCode >= 500 && resp.StatusCode < 600 && c.isIdempotent(req.Method) {
			c.circuitBreaker.recordFailure()
			serverErrorRetries++
			if serverErrorRetries > Max5xxRetries {
				return resp, nil
			}

			_ = resp.Body.Close()

			if c.debug {
				fmt.Fprintf(os.Stderr, "  server error, retry %d/%d in %v\n", serverErrorRetries, Max5xxRetries, ServerErrorRetryDelay)
			}

			if err := c.replayRequestBody(req); err != nil {
				return nil, fmt.Errorf("failed to replay request body: %w", err)
			}

			if err := c.sleep(ctx, ServerErrorRetryDelay); err != nil {
				return nil, err
			}
			continue
		}

		// Record success to reset circuit breaker
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			c.circuitBreaker.recordSuccess()
		}

		return resp, nil
	}
}

// parseRetryAfter parses the Retry-After header and returns the delay.
// Falls back to exponential backoff if header is missing or invalid.
func (c *Client) parseRetryAfter(header string, attempt int) time.Duration {
	if header != "" {
		if seconds, err := strconv.Atoi(header); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	// Exponential backoff: 1s, 2s, 4s, ...
	return RateLimitBaseDelay * (1 << (attempt - 1))
}

// isIdempotent returns true if the HTTP method is idempotent.
func (c *Client) isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodPut:
		return true
	default:
		return false
	}
}

// replayRequestBody resets the request body for retry using GetBody if available.
func (c *Client) replayRequestBody(req *http.Request) error {
	if req.GetBody == nil {
		return nil
	}
	body, err := req.GetBody()
	if err != nil {
		return err
	}
	req.Body = body
	return nil
}

// sleep waits for the given duration, respecting context cancellation.
func (c *Client) sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}

func (c *Client) Post(ctx context.Context, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	var bodyData []byte
	if body != nil {
		var err error
		bodyData, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyData)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	// Set GetBody for retry support
	if bodyData != nil {
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyData)), nil
		}
	}
	return c.Do(ctx, req)
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req)
}
