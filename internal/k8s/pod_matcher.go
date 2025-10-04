package k8s

import "regexp"

// MatchDefaultPod finds the first pod matching the given regex pattern.
// Returns the index of the matching pod and true if found, or -1 and false if no match.
// If pattern is nil, selects the first pod (index 0) if any pods exist.
func MatchDefaultPod(pods []Pod, pattern *regexp.Regexp) (int, bool) {
	if len(pods) == 0 {
		return -1, false
	}

	// If no pattern, select first pod
	if pattern == nil {
		return 0, true
	}

	// Find first matching pod
	for i, pod := range pods {
		if pattern.MatchString(pod.Name) {
			return i, true
		}
	}

	// No match found
	return -1, false
}
