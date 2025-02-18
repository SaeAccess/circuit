package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/podman"
)

func init() {
	client.RegisterElementMaker(&containerElementMaker{
		client.NewBaseElementMaker("container", ContainerType),
	})
}

var ContainerType = reflect.TypeOf((*podman.Container)(nil)).Elem()

// implementation for a specific maker
type containerElementMaker struct {
	client.BaseElementMaker
}
