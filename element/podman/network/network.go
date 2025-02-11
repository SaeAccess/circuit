package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"runtime"

	"github.com/gocircuit/circuit/anchor"
	c "github.com/gocircuit/circuit/client/podman"
	"github.com/gocircuit/circuit/element"
	"github.com/gocircuit/circuit/element/podman"
	"github.com/gocircuit/circuit/use/circuit"
)

// Network represents the network element
type Network interface {
	c.Network
	X() circuit.X
}

type network struct {
	name string
	id   string
}

func init() {
	anchor.RegisterElement(anchor.Network, ef, yf)
}

func MakeNetwork(opts c.NetworkCreateOptions) (Network, error) {
	// Check if podman enabled on this server
	// TODO make this a capability of the server when it joins the cluster.
	if podman.Path == "" {
		return nil, errors.New("podman not installed on this server")
	}

	// determine name
	opts.Name = element.ElementName(opts.Name)

	// Create a new exec.Cmd
	args := opts.CmdLine()
	log.Printf("cmd line: %s %v", podman.Path, args)
	r, err := exec.Command(podman.Path, args...).Output()
	if err != nil {
		log.Printf("error running command: %s %v - error:%v", podman.Path, args, err)
		return nil, err
	}

	// Create a new container
	netw := &network{
		name: opts.Name,
		id:   string(r),
	}

	// GC...
	runtime.SetFinalizer(netw,
		func(c *network) {
			exec.Command(podman.Path, "rm", c.name).Run()
		},
	)

	return netw, nil
}

func (n *network) Connect(opts *c.NetworkConnectOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(n.name)...).Run()
}

func (n *network) Disconnect(opts *c.NetworkDisconnectOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(n.name)...).Run()
}

func (n *network) Exists() error {
	return exec.Command(podman.Path, "network", "exists", n.name).Run()
}

func (n *network) Inspect() (*c.InspectNetworkSettings, error) {
	b, err := exec.Command(podman.Path, "network", "inspect", n.name).Output()
	if err != nil {
		return nil, err
	}

	data, err := element.ParseJSONArrayFirst[*c.InspectNetworkSettings](b)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n *network) PeekBytes() []byte {
	v, err := n.Inspect()
	if err != nil {
		// TODO Should the error be returned here as byte slise
		return []byte{}
	}

	b, _ := json.MarshalIndent(v, "", "\t")
	return b
}

func (n *network) Reload() error {
	return exec.Command(podman.Path, "network", "reload", n.name).Run()
}

func (n *network) Scrub() {
	opts := c.NetworkRemoveOptions{Force: true}
	exec.Command(podman.Path, opts.CmdLine(n.name)...).Run()
}

func (n *network) Update(opts *c.NetworkUpdateOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(n.name)...).Run()
}

func (n *network) X() circuit.X {
	return circuit.Ref(XNetwork{n})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	opts, ok := arg.(c.NetworkCreateOptions)
	if !ok {
		return nil, fmt.Errorf("invalid argument to network element factory, arg=%T", arg)
	}

	n, err := MakeNetwork(opts)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func yf(x circuit.X) (any, error) {
	return YNetwork{x}, nil
}
