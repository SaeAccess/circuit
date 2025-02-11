package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/podman/container"
)

func init() {
	client.RegisterElementMaker(&containerElementMaker{
		client.NewBaseElementMaker("container", reflect.TypeOf(container.YContainer{})),
	})
}

var ContainerType = reflect.TypeOf(container.YContainer{})

// implementation for a specific maker
type containerElementMaker struct {
	client.BaseElementMaker
}

func (c *containerElementMaker) Get(v any) any {
	return v.(container.YContainer)
}
