package executor

import "errors"

var (
	// ErrNoPodMatch indicates no pods match the pattern
	ErrNoPodMatch = errors.New("no pods match pattern")

	// ErrMultiplePodMatch indicates multiple pods match the pattern
	ErrMultiplePodMatch = errors.New("multiple pods match pattern")

	// ErrPodExecFailed indicates kubectl exec command failed
	ErrPodExecFailed = errors.New("kubectl exec failed")
)
