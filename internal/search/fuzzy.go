package search

import (
	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/sahilm/fuzzy"
)

// Match represents a fuzzy search match result
type Match struct {
	Namespace    k8s.Namespace
	MatchIndices []int // Character positions that matched the query
}

// FuzzyMatch performs fuzzy search on namespace list and returns matches with highlighting indices
func FuzzyMatch(query string, namespaces []k8s.Namespace) []Match {
	// Empty query returns all namespaces
	if query == "" {
		matches := make([]Match, len(namespaces))
		for i, ns := range namespaces {
			matches[i] = Match{
				Namespace:    ns,
				MatchIndices: []int{},
			}
		}
		return matches
	}

	// Build string slice for fuzzy library
	names := make([]string, len(namespaces))
	for i, ns := range namespaces {
		names[i] = ns.Name
	}

	// Perform fuzzy search
	fuzzyResults := fuzzy.Find(query, names)

	// Convert fuzzy results to Match structs
	matches := make([]Match, len(fuzzyResults))
	for i, result := range fuzzyResults {
		matches[i] = Match{
			Namespace:    namespaces[result.Index],
			MatchIndices: result.MatchedIndexes,
		}
	}

	return matches
}
