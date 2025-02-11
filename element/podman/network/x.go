package network

import (
	c "github.com/gocircuit/circuit/client/podman"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XNetwork{})
}

// X
type XNetwork struct {
	Network
}

func (x XNetwork) Connect(opts *c.NetworkConnectOptions) error {
	return errors.Pack(x.Network.Connect(opts))
}

func (x XNetwork) Disconnect(opts *c.NetworkDisconnectOptions) error {
	return errors.Pack(x.Network.Disconnect(opts))
}

func (x XNetwork) Exists() error {
	return errors.Pack(x.Network.Exists())
}

func (x XNetwork) Inspect() (*c.InspectNetworkSettings, error) {
	d, err := x.Network.Inspect()
	return d, errors.Pack(err)
}

func (x XNetwork) PeekBytes() []byte {
	return x.Network.PeekBytes()
}

func (x XNetwork) Reload() error {
	return errors.Pack(x.Network.Reload())
}

func (x XNetwork) Scrub() {
	x.Network.Scrub()
}

func (x XNetwork) Update(opts *c.NetworkUpdateOptions) error {
	return errors.Pack(x.Network.Update(opts))
}

type YNetwork struct {
	X circuit.X
}

func (y YNetwork) Connect(opts *c.NetworkConnectOptions) error {
	r := y.X.Call("Connect", opts)
	return errors.Unpack(r[0])
}

func (y YNetwork) Disconnect(opts *c.NetworkDisconnectOptions) error {
	r := y.X.Call("Disconnect", opts)
	return errors.Unpack(r[0])
}

func (y YNetwork) Exists() error {
	r := y.X.Call("Exists")
	return errors.Unpack(r[0])
}

func (y YNetwork) Inspect() (*c.InspectNetworkSettings, error) {
	r := y.X.Call("Inspect")
	data, _ := r[0].(*c.InspectNetworkSettings)
	return data, errors.Unpack(r[1])

}

func (y YNetwork) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YNetwork) Reload() error {
	r := y.X.Call("Reload")
	return errors.Unpack(r[0])
}

func (y YNetwork) Scrub() {
	y.X.Call("Scrub")
}

func (y YNetwork) Update(opts *c.NetworkUpdateOptions) error {
	r := y.X.Call("Update", opts)
	return errors.Unpack(r[0])
}
