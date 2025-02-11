package volume

import (
	c "github.com/gocircuit/circuit/client/podman"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XVolume{})
}

// X
type XVolume struct {
	Volume
}

func (x XVolume) Exists() error {
	return errors.Pack(x.Volume.Exists())
}

func (x XVolume) Export(opts c.VolumeExportOptions) error {
	return errors.Pack(x.Volume.Export(opts))
}

func (x XVolume) Import(source string) error {
	return errors.Pack(x.Volume.Import(source))
}

func (x XVolume) Inspect() (*c.InspectVolumeData, error) {
	r, err := x.Volume.Inspect()
	return r, errors.Pack(err)
}

func (x XVolume) Mount() error {
	return errors.Pack(x.Volume.Mount())
}

func (x XVolume) PeekBytes() []byte {
	return x.Volume.PeekBytes()
}

func (x XVolume) Reload() error {
	return errors.Pack(x.Volume.Reload())
}

func (x XVolume) Scrub() {
	x.Volume.Scrub()
}

func (x XVolume) Unmount() error {
	return errors.Pack(x.Volume.Unmount())
}

type YVolume struct {
	X circuit.X
}

func (y YVolume) Exists() error {
	r := y.X.Call("Exists")
	return errors.Unpack(r[0])
}

func (y YVolume) Export(opts c.VolumeExportOptions) error {
	r := y.X.Call("Export", opts)
	return errors.Unpack(r[0])
}

func (y YVolume) Import(source string) error {
	r := y.X.Call("Import", source)
	return errors.Unpack(r[0])
}

func (y YVolume) Inspect() (*c.InspectVolumeData, error) {
	r := y.X.Call("Inspect")
	data, _ := r[0].(*c.InspectVolumeData)
	return data, errors.Unpack(r[1])
}

func (y YVolume) Mount() error {
	r := y.X.Call("Mount")
	return errors.Unpack(r[0])
}

func (y YVolume) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YVolume) Reload() error {
	r := y.X.Call("Reload")
	return errors.Unpack(r[0])
}
func (y YVolume) Scrub() {
	y.X.Call("Scrub")
}

func (y YVolume) Unmount() error {
	r := y.X.Call("Unmount")
	return errors.Unpack(r[0])
}
