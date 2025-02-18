package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/docker"
)

func init() {
	client.RegisterElementMaker(&dockerElementMaker{
		client.NewBaseElementMaker("docker", DockerType),
	})
}

var DockerType = reflect.TypeOf((*docker.Container)(nil)).Elem()

// implementation for a specific maker
type dockerElementMaker struct {
	client.BaseElementMaker
}
