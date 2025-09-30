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

// Namespace represents a Kubernetes namespace
type Namespace struct {
	Name       string
	IsFavorite bool
}

// NamespaceList represents the JSON response from kubectl get namespaces
type NamespaceList struct {
	Items []NamespaceItem `json:"items"`
}

// NamespaceItem represents a single namespace in kubectl JSON output
type NamespaceItem struct {
	Metadata NamespaceMetadata `json:"metadata"`
}

// NamespaceMetadata contains namespace metadata
type NamespaceMetadata struct {
	Name string `json:"name"`
}

// Pod represents a Kubernetes pod (placeholder for future stories)
type Pod struct {
	Name   string
	Status string
}
