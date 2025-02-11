package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/docker"
)

func init() {
	client.RegisterElementMaker(&dockerElementMaker{
		client.NewBaseElementMaker("docker", reflect.TypeOf(docker.YContainer{})),
	})
}

var DockerType = reflect.TypeOf(docker.YContainer{})

// implementation for a specific maker
type dockerElementMaker struct {
	client.BaseElementMaker
}

func (c *dockerElementMaker) Get(v any) any {
	return v.(docker.YContainer)
}
