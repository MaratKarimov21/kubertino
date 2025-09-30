package k8s

// KubeConfig represents the structure of a kubeconfig file
type KubeConfig struct {
	Contexts       []KubeContext `yaml:"contexts"`
	CurrentContext string        `yaml:"current-context"`
	Clusters       []Cluster     `yaml:"clusters,omitempty"`
}

// KubeContext represents a context entry in kubeconfig
type KubeContext struct {
	Name    string         `yaml:"name"`
	Context ContextDetails `yaml:"context"`
}

// ContextDetails contains the cluster and namespace for a context
type ContextDetails struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace,omitempty"`
}

// Cluster represents a cluster entry in kubeconfig
type Cluster struct {
	Name    string        `yaml:"name"`
	Cluster ClusterConfig `yaml:"cluster"`
}

// ClusterConfig contains cluster connection details
type ClusterConfig struct {
	Server string `yaml:"server"`
}

// Pod represents a Kubernetes pod (placeholder for future stories)
type Pod struct {
	Name   string
	Status string
}
