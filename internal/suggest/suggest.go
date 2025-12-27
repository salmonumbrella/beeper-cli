// internal/suggest/suggest.go
package suggest

import (
	"fmt"
	"sort"
	"strings"
)

// Match represents a suggested match with its score
type Match struct {
	Value string
	Label string
	Score int
}

// FindSimilar finds items similar to the query using simple substring matching.
// Returns up to maxResults matches, sorted by relevance.
func FindSimilar(query string, items []Match, maxResults int) []Match {
	if maxResults <= 0 {
		return nil
	}
	if query == "" || len(items) == 0 {
		return nil
	}

	query = strings.ToLower(query)
	var matches []Match

	for _, item := range items {
		score := calculateScore(query, strings.ToLower(item.Value), strings.ToLower(item.Label))
		if score > 0 {
			item.Score = score
			matches = append(matches, item)
		}
	}

	// Sort by score descending, with stable tiebreaker by value
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].Value < matches[j].Value
	})

	if len(matches) > maxResults {
		matches = matches[:maxResults]
	}

	return matches
}

// calculateScore returns a relevance score (higher = better match)
func calculateScore(query, value, label string) int {
	score := 0

	// Exact match - highest priority
	if value == query {
		return 1000
	}

	// Starts with query
	if strings.HasPrefix(value, query) {
		score += 50
	} else if strings.Contains(value, query) {
		// Contains query (but doesn't start with it)
		score += 30
	}

	// Label contains query
	if strings.Contains(label, query) {
		score += 20
	}

	return score
}

// FormatSuggestions formats matches for display
func FormatSuggestions(matches []Match) string {
	if len(matches) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\nDid you mean one of these?\n")
	for _, m := range matches {
		if m.Label != "" {
			sb.WriteString(fmt.Sprintf("  %s  %s\n", m.Value, m.Label))
		} else {
			sb.WriteString(fmt.Sprintf("  %s\n", m.Value))
		}
	}
	return sb.String()
}
