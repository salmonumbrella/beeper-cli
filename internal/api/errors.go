package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// ValidationError represents a single validation error from the API
type ValidationError struct {
	Code     string   `json:"code"`
	Expected string   `json:"expected,omitempty"`
	Received string   `json:"received,omitempty"`
	Path     []string `json:"path,omitempty"`
	Message  string   `json:"message"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("API error: %d", e.StatusCode)
}

// parseValidationErrors attempts to parse array of validation errors
func parseValidationErrors(body []byte) (string, bool) {
	var validationErrors []ValidationError
	if err := json.Unmarshal(body, &validationErrors); err != nil {
		return "", false
	}
	if len(validationErrors) == 0 {
		return "", false
	}

	// Format validation errors nicely
	var messages []string
	for _, ve := range validationErrors {
		path := strings.Join(ve.Path, ".")
		if path != "" {
			messages = append(messages, fmt.Sprintf("%s: %s", path, ve.Message))
		} else {
			messages = append(messages, ve.Message)
		}
	}
	return strings.Join(messages, "; "), true
}

func ParseError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	apiErr := &APIError{StatusCode: resp.StatusCode}

	body, err := io.ReadAll(resp.Body)
	if err == nil && len(body) > 0 {
		// Try parsing as validation errors array first
		if msg, ok := parseValidationErrors(body); ok {
			apiErr.Message = msg
			return apiErr
		}
		_ = json.Unmarshal(body, apiErr) // Ignore unmarshal errors, we have defaults
	}

	return apiErr
}

func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// UserFriendlyError wraps an error with a user-friendly message
func UserFriendlyError(err error) error {
	if err == nil {
		return nil
	}

	// Connection refused - Beeper not running
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		if strings.Contains(netErr.Error(), "connection refused") {
			return fmt.Errorf("Beeper Desktop not running or API disabled. Start Beeper and enable Developer API in Settings") //nolint:staticcheck // User-facing error message
		}
	}

	// Check for connection refused in error string as fallback
	if strings.Contains(err.Error(), "connection refused") {
		return fmt.Errorf("Beeper Desktop not running or API disabled. Start Beeper and enable Developer API in Settings") //nolint:staticcheck // User-facing error message
	}

	return err
}

// ParseErrorWithContext returns user-friendly error from HTTP response
func ParseErrorWithContext(resp *http.Response, context string) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	apiErr := &APIError{StatusCode: resp.StatusCode}

	body, err := io.ReadAll(resp.Body)
	if err == nil && len(body) > 0 {
		// Try parsing as validation errors array first
		if msg, ok := parseValidationErrors(body); ok {
			apiErr.Message = msg
		} else {
			_ = json.Unmarshal(body, apiErr) // Ignore unmarshal errors, we have defaults
		}
	}

	// Handle specific status codes with friendly messages
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("Invalid or expired token. Run: beeper auth add") //nolint:staticcheck // User-facing error message
	case http.StatusNotFound:
		if context != "" {
			return fmt.Errorf("%s not found", context)
		}
		if apiErr.Message != "" {
			return fmt.Errorf("Not found: %s", apiErr.Message) //nolint:staticcheck // User-facing error message
		}
		return fmt.Errorf("Not found") //nolint:staticcheck // User-facing error message
	case http.StatusBadRequest, http.StatusInternalServerError:
		// Validation errors typically come as 400 or 500
		if apiErr.Message != "" {
			return fmt.Errorf("invalid request: %s", apiErr.Message)
		}
		return apiErr
	default:
		// Show Beeper's error message directly for other API errors
		if apiErr.Message != "" {
			return fmt.Errorf("%s", apiErr.Message)
		}
		return apiErr
	}
}
