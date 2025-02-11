package container

import (
	"io"

	c "github.com/gocircuit/circuit/client/podman"
	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XContainer{})
}

// XContainer is a circuit container that wraps a podman container skeleton impl.
type XContainer struct {
	Container
}

// Checkpoint the container
func (x XContainer) CheckPoint(opts *c.ContainerCheckpointOptions) error {
	return errors.Pack(x.Container.CheckPoint(opts))
}

// Clone an existing container
// func (x XContainer) Clone(opts *c.ContainerCloneOptions) (c.Container, error) {
// 	c, err := x.Container.Clone(opts)
// 	return c, errors.Pack(err)
// }

// Exec runs a process in the container
func (x XContainer) Exec(opts *c.ContainerExecOptions) ([]byte, error) {
	b, err := x.Container.Exec(opts)
	return b, errors.Pack(err)
}

// Inspect the configuration of the container
func (x XContainer) Inspect() (*c.InspectContainerData, error) {
	s, err := x.Container.Inspect()
	return s, errors.Pack(err)
}

func (x XContainer) IsDone() bool {
	return x.Container.IsDone()
}

// Pause the processes in the container
func (x XContainer) Pause() error {
	return errors.Pack(x.Container.Pause())
}

// Stats()
// Peek at the container's resource ussage statistics
func (x XContainer) Peek() (*c.InspectContainerData, error) {
	c, err := x.Container.Peek()
	return c, errors.Pack(err)
}

func (x XContainer) PeekBytes() []byte {
	return x.Container.PeekBytes()
}

// Get port mappings
func (x XContainer) Ports() []string {
	return x.Container.Ports()
}

// Restore the container from a checkpoint
func (x XContainer) Restore(opts *c.ContainerRestoreOptions) error {
	return errors.Pack(x.Container.Restore(opts))
}

// RunLabel runs the command specified in the label
func (x XContainer) RunLabel() error {
	return errors.Pack(x.Container.RunLabel())
}

// Scrub removes the container
func (x XContainer) Scrub() {
	x.Container.Scrub()
}

// Start the container
func (x XContainer) Start() error {
	return errors.Pack(x.Container.Start())
}

// Stop the container
func (x XContainer) Stop(opts *c.ContainerStopOpts) error {
	return errors.Pack(x.Container.Stop(opts))
}

// Unpause the processes in the container
func (x XContainer) Unpause() error {
	return errors.Pack(x.Container.Unpause())
}

// Kill the container with the specified signal
func (x XContainer) Signal(sig string) error {
	return errors.Pack(x.Container.Signal(sig))
}

// Wait for the container to exit
func (x XContainer) Wait() (*c.InspectContainerData, error) {
	s, err := x.Container.Wait()
	return s, errors.Pack(err)
}

func (x XContainer) Stdin() circuit.X {
	return xio.NewXWriteCloser(x.Container.Stdin())
}
func (x XContainer) Stdout() circuit.X {
	return xio.NewXReadCloser(x.Container.Stdout())

}
func (x XContainer) Stderr() circuit.X {
	return xio.NewXReadCloser(x.Container.Stderr())
}

// YContainer is a circuit container that wraps a podman container stub/proxy impl.
type YContainer struct {
	X circuit.X
}

func (y YContainer) CheckPoint(opts *c.ContainerCheckpointOptions) error {
	r := y.X.Call("CheckPoint", opts)
	return errors.Unpack(r[0])
}

// func (y YContainer) Clone() (con c.Container, err error) {
// 	r := y.X.Call("Clone")
// 	con, _ = r[0].(c.Container)
// 	return con, errors.Unpack(r[1])
// }

func (y YContainer) Exec(opts *c.ContainerExecOptions) ([]byte, error) {
	r := y.X.Call("Exec", opts)
	b, _ := r[0].([]byte)
	return b, errors.Unpack(r[1])
}

func (y YContainer) Inspect() (stat *c.InspectContainerData, err error) {
	r := y.X.Call("Inspect")
	stat, _ = r[0].(*c.InspectContainerData)
	return stat, errors.Unpack(r[1])
}

func (y YContainer) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YContainer) Pause() error {
	r := y.X.Call("Pause")
	return errors.Unpack(r[0])
}

func (y YContainer) Peek() (stat *c.InspectContainerData, err error) {
	r := y.X.Call("Peek")
	stat, _ = r[0].(*c.InspectContainerData)
	return stat, errors.Unpack(r[1])
}

func (y YContainer) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YContainer) Ports() []string {
	return y.X.Call("Ports")[0].([]string)
}

func (y YContainer) Restore(opts *c.ContainerRestoreOptions) error {
	r := y.X.Call("Restore")
	return errors.Unpack(r[0])
}

func (y YContainer) RunLabel() error {
	r := y.X.Call("RunLabel")
	return errors.Unpack(r[0])
}

func (y YContainer) Scrub() {
	y.X.Call("Scrub")
}

func (y YContainer) Start() error {
	r := y.X.Call("Start")
	return errors.Unpack(r[0])
}

func (y YContainer) Stop(opts *c.ContainerStopOpts) error {
	r := y.X.Call("Stop")
	return errors.Unpack(r[0])
}

func (y YContainer) Unpause() error {
	r := y.X.Call("Unpause")
	return errors.Unpack(r[0])
}

func (y YContainer) Signal(sig string) error {
	r := y.X.Call("Signal", sig)
	return errors.Unpack(r[0])
}

func (y YContainer) Wait() (stat *c.InspectContainerData, err error) {
	r := y.X.Call("Wait")
	stat, _ = r[0].(*c.InspectContainerData)
	return stat, errors.Unpack(r[1])
}

func (y YContainer) Stdin() io.WriteCloser {
	return xio.NewYWriteCloser(y.X.Call("Stdin")[0])
}

func (y YContainer) Stdout() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stdout")[0])
}

func (y YContainer) Stderr() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stderr")[0])
}
