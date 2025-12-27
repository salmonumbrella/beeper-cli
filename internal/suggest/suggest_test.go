// internal/suggest/suggest_test.go
package suggest

import (
	"strings"
	"testing"
)

func TestFindSimilar(t *testing.T) {
	items := []Match{
		{Value: "chats", Label: "List chats"},
		{Value: "chat", Label: "Chat command"},
		{Value: "messages", Label: "List messages"},
	}

	matches := FindSimilar("chat", items, 3)

	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}

	// First match should be exact
	if matches[0].Value != "chat" {
		t.Errorf("first match should be 'chat', got %q", matches[0].Value)
	}
}

func TestFindSimilarEmptyQuery(t *testing.T) {
	items := []Match{{Value: "test", Label: "Test"}}
	matches := FindSimilar("", items, 3)
	if matches != nil {
		t.Error("expected nil for empty query")
	}
}

func TestFindSimilarEmptyItems(t *testing.T) {
	matches := FindSimilar("test", []Match{}, 3)
	if matches != nil {
		t.Error("expected nil for empty items")
	}

	matches = FindSimilar("test", nil, 3)
	if matches != nil {
		t.Error("expected nil for nil items")
	}
}

func TestFindSimilarMaxResults(t *testing.T) {
	items := []Match{
		{Value: "abc", Label: ""},
		{Value: "abcd", Label: ""},
		{Value: "abcde", Label: ""},
	}

	// Test maxResults = 0
	matches := FindSimilar("abc", items, 0)
	if matches != nil {
		t.Error("expected nil for maxResults = 0")
	}

	// Test maxResults < 0
	matches = FindSimilar("abc", items, -1)
	if matches != nil {
		t.Error("expected nil for maxResults < 0")
	}

	// Test maxResults = 1
	matches = FindSimilar("abc", items, 1)
	if len(matches) != 1 {
		t.Errorf("expected 1 match, got %d", len(matches))
	}

	// Test maxResults = 2
	matches = FindSimilar("abc", items, 2)
	if len(matches) != 2 {
		t.Errorf("expected 2 matches, got %d", len(matches))
	}
}

func TestFindSimilarPrefixVsContains(t *testing.T) {
	items := []Match{
		{Value: "xyzabc", Label: ""}, // contains but doesn't start with
		{Value: "abcxyz", Label: ""}, // starts with abc
	}

	matches := FindSimilar("abc", items, 2)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	// Prefix match should score higher than contains match
	if matches[0].Value != "abcxyz" {
		t.Errorf("prefix match should be first, got %q", matches[0].Value)
	}
	if matches[1].Value != "xyzabc" {
		t.Errorf("contains match should be second, got %q", matches[1].Value)
	}

	// Verify scores: prefix gets 50, contains gets 30
	if matches[0].Score != 50 {
		t.Errorf("prefix match should have score 50, got %d", matches[0].Score)
	}
	if matches[1].Score != 30 {
		t.Errorf("contains match should have score 30, got %d", matches[1].Score)
	}
}

func TestFormatSuggestionsEmpty(t *testing.T) {
	result := FormatSuggestions(nil)
	if result != "" {
		t.Errorf("expected empty string for nil matches, got %q", result)
	}

	result = FormatSuggestions([]Match{})
	if result != "" {
		t.Errorf("expected empty string for empty matches, got %q", result)
	}
}

func TestFormatSuggestionsWithLabels(t *testing.T) {
	matches := []Match{
		{Value: "chat", Label: "Chat command"},
		{Value: "chats", Label: ""},
	}

	result := FormatSuggestions(matches)

	// Should contain header
	if !strings.Contains(result, "Did you mean") {
		t.Error("expected header in result")
	}

	// Should contain value with label
	if !strings.Contains(result, "chat") || !strings.Contains(result, "Chat command") {
		t.Error("expected 'chat' with its label in result")
	}

	// Should contain value without label
	if !strings.Contains(result, "chats") {
		t.Error("expected 'chats' in result")
	}
}
