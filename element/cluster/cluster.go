package cluster

import (
	"fmt"
	"sync"

	"github.com/gocircuit/circuit/anchor"
	cl "github.com/gocircuit/circuit/client/cluster"
	"github.com/gocircuit/circuit/use/circuit"
)

type Cluster interface {
	cl.Cluster
	X() circuit.X
}

type cluster struct {
	config *cl.ClusterConfig

	// circuit host address
	circuits []string

	//
}

func init() {
	anchor.RegisterElement("cluster", ef, yf)
}

func createCluster(config *cl.ClusterConfig) (_ Cluster, err error) {
	c := &cluster{
		config:   config,
		circuits: []string{},
	}
	// create and start circuit daemons on each of the host
	var wg sync.WaitGroup
	wg.Add(len(config.Hosts))
	// for _, h := range config.Hosts {
	// 	// config circuit process and start it
	// 	// first locate ciricuit
	// 	exe, err := exec.LookPath("circuit")
	// 	if err != nil {
	// 		return nil, errors.Wrapf(err, "circuit daemon not found")
	// 	}

	// 	cmd := proc.Cmd{
	// 		Path: exe,
	// 		Args: "",
	// 	}

	// }

	return c, nil
}

// AddHost adds a new host to the cluster
func (c *cluster) AddHost(hc *cl.HostConfig) error {
	return nil
}

// Clone clones the cluster
func (c *cluster) Clone() error {
	return nil
}

func (c *cluster) CreateCluster(*cl.ClusterConfig) error {
	return nil
}

// Exists checks if host exist in cluster
func (c *cluster) Exists(hostName string) bool {
	return false
}

// List the hosts in the cluster
func (c *cluster) Hosts() []string {
	return nil
}

// Inspect the cluster configuration
func (c *cluster) Inspect() (*cl.InspectClusterConfig, error) {
	return nil, nil
}

func (c *cluster) InspectHost(hostName string) (*cl.HostInfo, error) {
	return nil, nil
}

// Join merges the networks of this cluster with argument cluster
func (c *cluster) Join(cluster string) error {
	return nil
}

// RemoveHost terminates and removes the specified host in the cluster
func (c *cluster) RemoveHost(hostName string) error {
	return nil
}

// ResolveHost
func (c *cluster) ResolveHost(hostName string) (*cl.HostConfig, error) {
	return nil, nil
}

// Scrub shutdown and remove the cluster element
func (c *cluster) Scrub() {

}

// Shutdown the entire cluster, stops all hosts
func (c *cluster) Shutdown() error {
	return nil
}

// Signal all the hosts in the cluster
func (c *cluster) Signal(sig string) error {
	return nil
}

// Stack provides the runtime statck trace for the requested server at host
func (c *cluster) Stack(hostName string) {

}

// Start a stopped host in the cluster
func (c *cluster) Start(hostName string) error {
	return nil
}

// Stop stops the request server host in the cluster
func (c *cluster) Stop(hostName string) error {
	return nil
}

// Top displays the running processes
func (c *cluster) Top(hostName string) {

}

func (c *cluster) X() circuit.X {
	return circuit.Ref(XCluster{c})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	opts, ok := arg.(c.ClusterConfig)
	if !ok {
		return nil, fmt.Errorf("invalid argument to network element factory, arg=%T", arg)
	}

	n, err := createCluster(&opts)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func yf(x circuit.X) (any, error) {
	return YCluster{x}, nil
}
