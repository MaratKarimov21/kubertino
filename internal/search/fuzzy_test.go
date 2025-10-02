package search

import (
	"testing"

	"github.com/maratkarimov/kubertino/internal/k8s"
	"github.com/stretchr/testify/assert"
)

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		namespaces []k8s.Namespace
		wantCount  int
		wantFirst  string
	}{
		{
			name:  "matches non-contiguous chars",
			query: "kube",
			namespaces: []k8s.Namespace{
				{Name: "kube-system"},
				{Name: "default"},
				{Name: "kubertino-prod"},
			},
			wantCount: 2,
			wantFirst: "kube-system",
		},
		{
			name:  "empty query returns all",
			query: "",
			namespaces: []k8s.Namespace{
				{Name: "kube-system"},
				{Name: "default"},
			},
			wantCount: 2,
		},
		{
			name:  "no matches",
			query: "xyz",
			namespaces: []k8s.Namespace{
				{Name: "kube-system"},
				{Name: "default"},
			},
			wantCount: 0,
		},
		{
			name:  "case insensitive matching",
			query: "KUBE",
			namespaces: []k8s.Namespace{
				{Name: "kube-system"},
				{Name: "default"},
			},
			wantCount: 1,
			wantFirst: "kube-system",
		},
		{
			name:  "matches with favorites",
			query: "sys",
			namespaces: []k8s.Namespace{
				{Name: "kube-system", IsFavorite: true},
				{Name: "default", IsFavorite: false},
				{Name: "system-logs", IsFavorite: true},
			},
			wantCount: 2,
			// Don't assert specific order - fuzzy library determines ranking
		},
		{
			name:  "single character query",
			query: "d",
			namespaces: []k8s.Namespace{
				{Name: "default"},
				{Name: "kube-system"},
				{Name: "production"},
			},
			wantCount: 2, // default, production
		},
		{
			name:       "empty namespace list",
			query:      "test",
			namespaces: []k8s.Namespace{},
			wantCount:  0,
		},
		{
			name:  "special characters in namespace names",
			query: "app",
			namespaces: []k8s.Namespace{
				{Name: "my-app-prod"},
				{Name: "another-app"},
				{Name: "application"},
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := FuzzyMatch(tt.query, tt.namespaces)
			assert.Equal(t, tt.wantCount, len(matches), "unexpected number of matches")

			if tt.wantCount > 0 && tt.wantFirst != "" {
				assert.Equal(t, tt.wantFirst, matches[0].Namespace.Name, "unexpected first match")
			}

			// Verify that MatchIndices are provided for non-empty queries
			if tt.query != "" && len(matches) > 0 {
				// At least the first match should have indices
				assert.NotNil(t, matches[0].MatchIndices, "match indices should not be nil")
			}

			// For empty query, verify all namespaces are returned with empty indices
			if tt.query == "" && len(tt.namespaces) > 0 {
				for i, match := range matches {
					assert.Equal(t, tt.namespaces[i].Name, match.Namespace.Name)
					assert.Equal(t, 0, len(match.MatchIndices), "empty query should have no match indices")
				}
			}
		})
	}
}

func TestFuzzyMatch_PreservesFavoriteFlag(t *testing.T) {
	namespaces := []k8s.Namespace{
		{Name: "kube-system", IsFavorite: true},
		{Name: "default", IsFavorite: false},
	}

	matches := FuzzyMatch("kube", namespaces)

	assert.Equal(t, 1, len(matches))
	assert.True(t, matches[0].Namespace.IsFavorite, "favorite flag should be preserved")
}

func TestFuzzyMatch_MatchIndices(t *testing.T) {
	namespaces := []k8s.Namespace{
		{Name: "kube-system"},
	}

	matches := FuzzyMatch("ks", namespaces)

	assert.Equal(t, 1, len(matches))
	assert.NotNil(t, matches[0].MatchIndices, "match indices should not be nil")
	assert.Greater(t, len(matches[0].MatchIndices), 0, "should have match indices")
	// The fuzzy library should have found 'k' and 's' in "kube-system"
	assert.Contains(t, matches[0].MatchIndices, 0, "should match 'k' at position 0")
}
