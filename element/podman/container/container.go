package container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"

	"github.com/gocircuit/circuit/anchor"
	c "github.com/gocircuit/circuit/client/podman"
	"github.com/gocircuit/circuit/element"
	"github.com/gocircuit/circuit/element/podman"
	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/use/circuit"
)

type Container interface {
	c.Container
	X() circuit.X
}

type container struct {
	name   string
	id     string
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	exit   chan error
}

func init() {
	anchor.RegisterElement(anchor.Container, ef, yf)
}

// MakeContainer creates a new container element
func makeContainer(opts *c.ContainerCreateOptions) (Container, error) {
	// Check if podman enabled on this server
	// TODO make this a capability of the server when it joins the cluster.
	if podman.Path == "" {
		return nil, errors.New("podman not installed on this server")
	}

	// determine name
	opts.Name = element.ElementName(opts.Name)

	// Create a new exec.Cmd
	args := opts.CmdLine(opts.Name)
	log.Printf("cmd line: %s %v", podman.Path, args)
	cmd := exec.Command(podman.Path, args...)
	r, err := cmd.Output()
	if err != nil {
		log.Printf("error running command: %s %v - error:%v", podman.Path, args, err)
		return nil, err
	}

	// Create a new container
	con := &container{
		name: opts.Name,
		id:   string(r),
		exit: make(chan error, 1),
		cmd:  cmd,
	}

	// GC...
	runtime.SetFinalizer(con,
		func(c *container) {
			exec.Command(podman.Path, "rm", c.name).Run()
		},
	)

	return con, nil
}

// NOTE: this requires root to run and also that criu is installed on the host.
// TODO check user and if criu is installed from host config
func (con *container) CheckPoint(opts *c.ContainerCheckpointOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(con.name)...).Run()
}

// Exec runs a command in a running container.
func (con *container) Exec(opts *c.ContainerExecOptions) ([]byte, error) {
	return exec.Command(podman.Path, opts.CmdLine(con.name)...).Output()
}

func (con *container) Inspect() (*c.InspectContainerData, error) {
	args := append([]string{}, "container", "inspect", con.name)

	b, err := exec.Command(podman.Path, args...).Output()
	if err != nil {
		return nil, err
	}

	data, err := element.ParseJSONArrayFirst[*c.InspectContainerData](b)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *container) IsDone() bool {
	select {
	case <-con.exit:
		return true
	default:
		return false
	}
}

func (con *container) Pause() error {
	if err := exec.Command(podman.Path, "container", "pause", con.name).Run(); err != nil {
		return err
	}

	return nil
}

func (con *container) Peek() (*c.InspectContainerData, error) {
	return con.Inspect()
}

func (con *container) PeekBytes() []byte {
	v, err := con.Inspect()
	if err != nil {
		return []byte{}
	}

	b, _ := json.MarshalIndent(v, "", "\t")
	return b
}

func (con *container) Ports() []string {
	b, err := exec.Command(podman.Path, "container", "port", con.name).Output()
	if err != nil {
		return []string{}
	}

	// convert to []string
	data, err := element.ParseJSONArray[string](b)
	if err != nil {
		return []string{}
	}

	// TODO what does b look like with multiple ports.

	return data
}

func (con *container) Restore(opts *c.ContainerRestoreOptions) error {
	return exec.Command(podman.Path, opts.CmdLine(con.name)...).Run()
}

func (con *container) RunLabel() error {
	return nil
}

func (con *container) Scrub() {
	// TODO stop, then remove
	opts := c.ContainerRemoveOptions{
		Depend: true, // remove any dependent container
		Force:  true,
		Ignore: true,
		Volume: true,
	}
	exec.Command(podman.Path, opts.CmdLine(con.name)...).Run()
}

func (con *container) Signal(sig string) error {
	// signal can be integer string or string name
	return nil
}

func (con *container) Start() error {
	if con.id == "" {
		return errors.New("container not created")
	}

	con.cmd = exec.Command(podman.Path, "start", con.name)

	con.cmd.Stdin, con.stdin = interruptible.BufferPipe(element.StdBufferLen)
	con.stdout, con.cmd.Stdout = interruptible.BufferPipe(element.StdBufferLen)
	con.stderr, con.cmd.Stderr = interruptible.BufferPipe(element.StdBufferLen)
	if err := con.cmd.Start(); err != nil {
		return err
	}

	go func() {
		con.exit <- con.cmd.Wait()
		close(con.exit)
		con.cmd.Stdout.(io.Closer).Close()
		con.cmd.Stderr.(io.Closer).Close()
	}()

	return nil
}

func (con *container) Stderr() io.ReadCloser {
	return con.stderr
}

func (con *container) Stdin() io.WriteCloser {
	return con.stdin
}

func (con *container) Stdout() io.ReadCloser {
	return con.stdout
}

func (con *container) Stop(opts *c.ContainerStopOpts) error {
	if err := exec.Command(podman.Path, opts.CmdLine(con.name)...).Run(); err != nil {
		return err
	}

	return nil
}

func (con *container) Unpause() error {
	if err := exec.Command(podman.Path, "container", "unpause", con.name).Run(); err != nil {
		return err
	}

	return nil
}

func (con *container) Wait() (*c.InspectContainerData, error) {
	<-con.exit
	return con.Peek()
}

func (con *container) X() circuit.X {
	return circuit.Ref(XContainer{con})
}

// ef is the element factory for the container element
func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	opts, ok := arg.(c.ContainerCreateOptions)
	if !ok {
		return nil, fmt.Errorf("invalid argument, arg=%T", arg)
	}
	x, err := makeContainer(&opts)
	if err != nil {
		return nil, err
	}
	go func() {
		defer func() {
			recover()
		}()
		if opts.Scrub {
			defer t.Scrub()
		}
		x.Wait()
	}()
	return x, nil
}

// yf is the element factory for the container element
func yf(x circuit.X) (any, error) {
	return YContainer{x}, nil
}
