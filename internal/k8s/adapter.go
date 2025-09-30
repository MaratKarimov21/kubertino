package k8s

// KubeAdapter abstracts Kubernetes operations
type KubeAdapter interface {
	// GetContexts returns the list of available kubectl contexts
	GetContexts() ([]string, error)

	// GetNamespaces returns namespaces for a given context (future implementation)
	GetNamespaces(context string) ([]string, error)

	// GetPods returns pods for a given context and namespace (future implementation)
	GetPods(context, namespace string) ([]Pod, error)
}
