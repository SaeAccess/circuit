package volume

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

type Volume interface {
	c.Volume
	X() circuit.X
}

type volume struct {
	name string
	id   string
}

func init() {
	anchor.RegisterElement(anchor.Volume, ef, yf)
}

func MakeVolume(opts c.VolumeCreateOptions) (Volume, error) {
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
	vol := &volume{
		name: opts.Name,
		id:   string(r),
	}

	// GC...
	runtime.SetFinalizer(vol,
		func(c *volume) {
			exec.Command(podman.Path, "rm", c.name).Run()
		},
	)

	return vol, nil

}

func (v *volume) Exists() error {
	return exec.Command(podman.Path, "volume", "exists", v.name).Run()
}

func (v *volume) Export(opts c.VolumeExportOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(v.name)...).Run()
}

func (v *volume) Import(source string) error {
	args := []string{"volume", "import", v.name}
	if source != "" {
		args = append(args, source)
	}

	return exec.Command(podman.Path, args...).Run()
}

func (v *volume) Inspect() (*c.InspectVolumeData, error) {
	b, err := exec.Command(podman.Path, "volume", "inspect", v.name).Output()
	if err != nil {
		return nil, err
	}

	data, err := element.ParseJSONArrayFirst[*c.InspectVolumeData](b)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (v *volume) Mount() error {
	return exec.Command(podman.Path, "volume", "mount", v.name).Run()
}

func (v *volume) PeekBytes() []byte {
	vol, err := v.Inspect()
	if err != nil {
		return []byte{}
	}

	b, _ := json.MarshalIndent(vol, "", "\t")
	return b
}

func (v *volume) Reload() error {
	return exec.Command(podman.Path, "volume", "reload", v.name).Run()
}

func (v *volume) Scrub() {
	opts := c.VolumeRemoveOptions{Force: true}
	exec.Command(podman.Path, opts.CmdLine(v.name)...).Run()
}

func (v *volume) Unmount() error {
	return exec.Command(podman.Path, "volume", "unmount", v.name).Run()
}

func (v *volume) X() circuit.X {
	return circuit.Ref(XVolume{v})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	opts, ok := arg.(c.VolumeCreateOptions)
	if !ok {
		return nil, fmt.Errorf("invalid argument to volume element factory, arg=%T", arg)
	}

	v, err := MakeVolume(opts)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func yf(x circuit.X) (any, error) {
	return YVolume{x}, nil
}
