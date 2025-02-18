package cluster

// Cluster (maintains context for a circuit cluster)
type Cluster interface {
	// AddHost adds a new host to the cluster
	AddHost(hc *HostConfig) error

	// Clone clones the cluster
	Clone() error

	CreateCluster(*ClusterConfig) error

	// Exists checks if host exist in cluster
	Exists(hostName string) bool

	// List the hosts in the cluster
	Hosts() []string

	// Inspect the cluster configuration
	Inspect() (*InspectClusterConfig, error)

	InspectHost(hostName string) (*HostInfo, error)

	// Join merges the networks of this cluster with argument cluster
	Join(cluster string) error

	// RemoveHost terminates and removes the specified host in the cluster
	RemoveHost(hostName string) error

	// ResolveHost
	ResolveHost(hostName string) (*HostConfig, error)

	// Scrub shutdown and remove the cluster element
	Scrub()

	// Shutdown the entire cluster, stops all hosts
	Shutdown() error

	// Signal all the hosts in the cluster
	Signal(sig string) error

	// Stack provides the runtime statck trace for the requested server at host
	Stack(hostName string)

	// Start a stopped host in the cluster
	Start(hostName string) error

	// Stop stops the request server host in the cluster
	Stop(hostName string) error

	// Stats provide live stream of resouce usage statistics for a server in the cluster
	//Stats(name string)

	// Top displays the running processes
	Top(hostName string)
}

type ClusterConfig struct {
	Name  string       `json:"name,omitempty"`
	Hosts []HostConfig `json:"hosts,omitempty"`
}

type HostConfig struct {
	// Name of host
	Name string

	// Address of the host
	Addr string `json:"addr,omitempty"`

	// Host network interface
	If string `json:"if,omitempty"`

	//
	Hmac string `json:"hmac,omitempty"`

	// Export to dns flag
}

type InspectClusterConfig struct {
}

type HostInfo struct {
}
