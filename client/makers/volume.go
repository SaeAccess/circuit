package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/podman/volume"
)

func init() {
	client.RegisterElementMaker(&volumeElementMaker{
		client.NewBaseElementMaker("volume", reflect.TypeOf(volume.YVolume{})),
	})
}

// implementation for a specific maker
type volumeElementMaker struct {
	client.BaseElementMaker
}

func (c *volumeElementMaker) Get(v any) any {
	return v.(volume.YVolume)
}
