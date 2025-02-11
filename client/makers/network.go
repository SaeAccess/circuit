package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/podman/network"
)

func init() {
	client.RegisterElementMaker(&networkElementMaker{
		client.NewBaseElementMaker("network", reflect.TypeOf(network.YNetwork{})),
	})
}

var NetworkType = reflect.TypeOf(network.YNetwork{})

// implementation for a specific maker
type networkElementMaker struct {
	client.BaseElementMaker
}

func (c *networkElementMaker) Get(v any) any {
	return v.(network.YNetwork)
}
