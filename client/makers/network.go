package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/podman"
)

func init() {
	client.RegisterElementMaker(&networkElementMaker{
		client.NewBaseElementMaker("network", NetworkType),
	})
}

var NetworkType = reflect.TypeOf((*podman.Network)(nil)).Elem()

// implementation for a specific maker
type networkElementMaker struct {
	client.BaseElementMaker
}
