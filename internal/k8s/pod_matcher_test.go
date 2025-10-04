package k8s

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchDefaultPod(t *testing.T) {
	tests := []struct {
		name      string
		pods      []Pod
		pattern   string
		wantIndex int
		wantFound bool
	}{
		{
			name: "single match",
			pods: []Pod{
				{Name: "web-frontend"},
				{Name: "api-server"},
				{Name: "worker"},
			},
			pattern:   "^api-.*",
			wantIndex: 1,
			wantFound: true,
		},
		{
			name: "multiple matches - first selected",
			pods: []Pod{
				{Name: "api-server-1"},
				{Name: "api-server-2"},
				{Name: "worker"},
			},
			pattern:   "^api-.*",
			wantIndex: 0,
			wantFound: true,
		},
		{
			name: "no match",
			pods: []Pod{
				{Name: "web-frontend"},
				{Name: "worker"},
			},
			pattern:   "^api-.*",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name:      "empty pod list",
			pods:      []Pod{},
			pattern:   ".*",
			wantIndex: -1,
			wantFound: false,
		},
		{
			name: "nil pattern - match first",
			pods: []Pod{
				{Name: "pod-1"},
				{Name: "pod-2"},
			},
			pattern:   "", // Will be nil after compilation check
			wantIndex: 0,
			wantFound: true,
		},
		{
			name: "match all pattern",
			pods: []Pod{
				{Name: "pod-1"},
				{Name: "pod-2"},
			},
			pattern:   ".*",
			wantIndex: 0,
			wantFound: true,
		},
		{
			name: "complex regex pattern",
			pods: []Pod{
				{Name: "backend-api-v1"},
				{Name: "frontend-ui-v1"},
				{Name: "backend-api-v2"},
			},
			pattern:   "^backend-.*-v2$",
			wantIndex: 2,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var compiled *regexp.Regexp
			if tt.pattern != "" {
				compiled = regexp.MustCompile(tt.pattern)
			}

			idx, found := MatchDefaultPod(tt.pods, compiled)

			assert.Equal(t, tt.wantIndex, idx)
			assert.Equal(t, tt.wantFound, found)
		})
	}
}
