package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/podman"
)

func init() {
	client.RegisterElementMaker(&podElementMaker{
		client.NewBaseElementMaker("pod", PodType),
	})
}

var PodType = reflect.TypeOf((*podman.Pod)(nil)).Elem()

// implementation for a specific maker
type podElementMaker struct {
	client.BaseElementMaker
}
