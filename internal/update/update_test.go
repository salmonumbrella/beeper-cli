// internal/update/update_test.go
package update

import (
	"context"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.0.0", "v1.0.0"},
		{"v1.0.0", "v1.0.0"},
		{"0.1.0", "v0.1.0"},
	}

	for _, tt := range tests {
		result := normalizeVersion(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeVersion(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestCheckForUpdateSkipsDevBuilds(t *testing.T) {
	ctx := context.Background()

	result := CheckForUpdate(ctx, "dev")
	if result != nil {
		t.Error("expected nil for dev builds")
	}

	result = CheckForUpdate(ctx, "")
	if result != nil {
		t.Error("expected nil for empty version")
	}
}
