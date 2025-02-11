package pod

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

type Pod interface {
	c.Pod
	X() circuit.X
}

type pod struct {
	name string
	id   string
}

func init() {
	anchor.RegisterElement(anchor.Pod, ef, yf)
}

func MakePod(opts c.PodCreateOptions) (Pod, error) {
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
		log.Printf("error running command, MakePod: %s %v - error:%v", podman.Path, args, err)
		return nil, err
	}

	// Create a new container
	p := &pod{
		name: opts.Name,
		id:   string(r),
	}

	// GC...
	runtime.SetFinalizer(p,
		func(c *pod) {
			exec.Command(podman.Path, "rm", c.name).Run()
		},
	)

	return p, nil
}

func (p *pod) Clone(opts *c.PodCloneOptions) (c.Pod, error) {
	b, err := exec.Command(podman.Path, opts.CmdLine()...).Output()
	if err != nil {
		return nil, err
	}

	// create a new pod anchor
	np := &pod{
		name: "", // TODO get name for new pod
		id:   string(b),
	}

	return np, nil
}

// TODO capture error (1 returned from command) and return pod not exist error
func (p *pod) Exists() error {
	return exec.Command(podman.Path, "pod", "exists", p.name).Run()
}

func (p *pod) Inspect() (*c.InspectPodData, error) {
	b, err := exec.Command(podman.Path, "pod", "inspect", p.name).Output()
	if err != nil {
		return nil, err
	}

	data, err := element.ParseJSONArrayFirst[*c.InspectPodData](b)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func (p *pod) Pause() error {
	return exec.Command(podman.Path, "pod", "pause", p.name).Run()
}

func (p *pod) PeekBytes() []byte {
	v, err := p.Inspect()
	if err != nil {
		return []byte{}
	}

	b, _ := json.MarshalIndent(v, "", "\t")
	return b
}

func (p *pod) Restart() error {
	return exec.Command(podman.Path, "pod", "restart", p.name).Run()
}

func (p *pod) Scrub() {
	// TODO should force removal and ignore when pod is missing
	args := []string{"pod", "rm", "--force", "--ignore", p.name}
	exec.Command(podman.Path, args...).Run()
}

func (p *pod) Signal(sig string) error {
	args := []string{"pod", "kill"}
	if sig != "" {
		args = append(args, "--signal", sig)
	}
	args = append(args, p.name)
	return exec.Command(podman.Path, args...).Run()
}

func (p *pod) Start(opts *c.PodStartOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(p.name)...).Run()
}

func (p *pod) Stop(opts *c.PodStopOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(p.name)...).Run()
}

func (p *pod) Unpause() error {
	return exec.Command(podman.Path, "pod", "unpause", p.name).Run()
}

func (p *pod) X() circuit.X {
	return circuit.Ref(XPod{p})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	opts, ok := arg.(c.PodCreateOptions)
	if !ok {
		return nil, fmt.Errorf("invalid argument to network element factory, arg=%T", arg)
	}

	n, err := MakePod(opts)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func yf(x circuit.X) (any, error) {
	return YPod{x}, nil
}
