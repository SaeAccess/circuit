package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/podman"
)

var VolumeType = reflect.TypeOf((*podman.Volume)(nil)).Elem()

func init() {
	client.RegisterElementMaker(&volumeElementMaker{
		client.NewBaseElementMaker("volume", VolumeType),
	})
}

// implementation for a specific maker
type volumeElementMaker struct {
	client.BaseElementMaker
}
