package k8s

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPods_ParseJSON(t *testing.T) {
	tests := []struct {
		name          string
		jsonOutput    string
		expectedPods  []Pod
		expectedError bool
	}{
		{
			name: "successful parse with multiple pods",
			jsonOutput: `{
				"items": [
					{
						"metadata": {"name": "pod-1"},
						"status": {"phase": "Running"}
					},
					{
						"metadata": {"name": "pod-2"},
						"status": {"phase": "Pending"}
					}
				]
			}`,
			expectedPods: []Pod{
				{Name: "pod-1", Status: "Running"},
				{Name: "pod-2", Status: "Pending"},
			},
			expectedError: false,
		},
		{
			name:          "empty pod list",
			jsonOutput:    `{"items": []}`,
			expectedPods:  []Pod{},
			expectedError: false,
		},
		{
			name: "single pod",
			jsonOutput: `{
				"items": [
					{
						"metadata": {"name": "single-pod"},
						"status": {"phase": "Failed"}
					}
				]
			}`,
			expectedPods: []Pod{
				{Name: "single-pod", Status: "Failed"},
			},
			expectedError: false,
		},
		{
			name: "all pod statuses",
			jsonOutput: `{
				"items": [
					{"metadata": {"name": "running-pod"}, "status": {"phase": "Running"}},
					{"metadata": {"name": "pending-pod"}, "status": {"phase": "Pending"}},
					{"metadata": {"name": "failed-pod"}, "status": {"phase": "Failed"}},
					{"metadata": {"name": "succeeded-pod"}, "status": {"phase": "Succeeded"}},
					{"metadata": {"name": "unknown-pod"}, "status": {"phase": "Unknown"}}
				]
			}`,
			expectedPods: []Pod{
				{Name: "running-pod", Status: "Running"},
				{Name: "pending-pod", Status: "Pending"},
				{Name: "failed-pod", Status: "Failed"},
				{Name: "succeeded-pod", Status: "Succeeded"},
				{Name: "unknown-pod", Status: "Unknown"},
			},
			expectedError: false,
		},
		{
			name:          "invalid JSON",
			jsonOutput:    `invalid json`,
			expectedPods:  nil,
			expectedError: true,
		},
		{
			name:          "malformed JSON structure",
			jsonOutput:    `{"items": "not an array"}`,
			expectedPods:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse JSON directly to simulate kubectl output parsing
			var response PodList
			err := json.Unmarshal([]byte(tt.jsonOutput), &response)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Convert to []Pod
			pods := make([]Pod, 0, len(response.Items))
			for _, item := range response.Items {
				pods = append(pods, Pod{
					Name:   item.Metadata.Name,
					Status: item.Status.Phase,
				})
			}

			assert.Equal(t, tt.expectedPods, pods)
		})
	}
}

func TestValidateNamespaceName(t *testing.T) {
	tests := []struct {
		name          string
		namespaceName string
		expectError   bool
	}{
		{
			name:          "valid lowercase",
			namespaceName: "default",
			expectError:   false,
		},
		{
			name:          "valid with hyphens",
			namespaceName: "kube-system",
			expectError:   false,
		},
		{
			name:          "valid with numbers",
			namespaceName: "namespace-123",
			expectError:   false,
		},
		{
			name:          "valid single character",
			namespaceName: "a",
			expectError:   false,
		},
		{
			name:          "invalid uppercase",
			namespaceName: "MyNamespace",
			expectError:   true,
		},
		{
			name:          "invalid underscore",
			namespaceName: "my_namespace",
			expectError:   true,
		},
		{
			name:          "invalid dot",
			namespaceName: "my.namespace",
			expectError:   true,
		},
		{
			name:          "invalid empty",
			namespaceName: "",
			expectError:   true,
		},
		{
			name:          "invalid starts with hyphen",
			namespaceName: "-namespace",
			expectError:   true,
		},
		{
			name:          "invalid ends with hyphen",
			namespaceName: "namespace-",
			expectError:   true,
		},
		{
			name:          "invalid special characters",
			namespaceName: "namespace@test",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNamespaceName(tt.namespaceName)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPodTypes(t *testing.T) {
	// Test Pod struct
	pod := Pod{
		Name:   "test-pod",
		Status: "Running",
	}
	assert.Equal(t, "test-pod", pod.Name)
	assert.Equal(t, "Running", pod.Status)

	// Test PodList JSON unmarshaling
	jsonStr := `{
		"items": [
			{
				"metadata": {"name": "pod-1"},
				"status": {"phase": "Running"}
			}
		]
	}`

	var podList PodList
	err := json.Unmarshal([]byte(jsonStr), &podList)
	require.NoError(t, err)
	assert.Len(t, podList.Items, 1)
	assert.Equal(t, "pod-1", podList.Items[0].Metadata.Name)
	assert.Equal(t, "Running", podList.Items[0].Status.Phase)
}
