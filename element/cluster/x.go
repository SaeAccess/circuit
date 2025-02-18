package cluster

import "github.com/gocircuit/circuit/use/circuit"

type XCluster struct {
	Cluster
}

func init() {
	circuit.RegisterValue(XCluster{})
}

type YCluster struct {
	X circuit.X
}
