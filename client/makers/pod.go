package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/podman/pod"
)

func init() {
	client.RegisterElementMaker(&podElementMaker{
		client.NewBaseElementMaker("pod", reflect.TypeOf(pod.YPod{})),
	})
}

var PodType = reflect.TypeOf(pod.YPod{})

// implementation for a specific maker
type podElementMaker struct {
	client.BaseElementMaker
}

func (c *podElementMaker) Get(v any) any {
	return v.(pod.YPod)
}
