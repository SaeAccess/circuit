package pod

import (
	"github.com/gocircuit/circuit/client/podman"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XPod{})
}

type XPod struct {
	Pod
}

// skeleton adapter
func (x XPod) Clone(opts *podman.PodCloneOptions) (podman.Pod, error) {
	p, err := x.Pod.Clone(opts)
	return p, errors.Pack(err)
}

func (x XPod) Exists() error {
	return errors.Pack(x.Pod.Exists())
}

func (x XPod) Inspect() (*podman.InspectPodData, error) {
	d, err := x.Pod.Inspect()
	return d, errors.Pack(err)
}

func (x XPod) Pause() error {
	return errors.Pack(x.Pod.Pause())
}

func (x XPod) PeekBytes() []byte {
	return x.Pod.PeekBytes()
}

func (x XPod) Restart() error {
	return errors.Pack(x.Pod.Restart())
}

func (x XPod) Scrub() {
	x.Pod.Scrub()
}

func (x XPod) Signal(sig string) error {
	return errors.Pack(x.Pod.Signal(sig))
}

func (x XPod) Start(opts *podman.PodStartOptions) error {
	return errors.Pack(x.Pod.Start(opts))
}

func (x XPod) Stop(opts *podman.PodStopOptions) error {
	return errors.Pack(x.Pod.Stop(opts))
}

func (x XPod) Unpause() error {
	return errors.Pack(x.Pod.Unpause())
}

// stub proxy
type YPod struct {
	X circuit.X
}

func (y YPod) Clone(opts *podman.PodCloneOptions) (podman.Pod, error) {
	r := y.X.Call("Clone", opts)
	p, _ := r[0].(podman.Pod)
	return p, errors.Unpack(r[1])
}

func (y YPod) Exists() error {
	r := y.X.Call("Exists")
	return errors.Unpack(r[0])
}

func (y YPod) Inspect() (*podman.InspectPodData, error) {
	r := y.X.Call("Inspect")
	d, _ := r[0].(*podman.InspectPodData)
	return d, errors.Unpack(r[1])
}

func (y YPod) Pause() error {
	r := y.X.Call("Pause")
	return errors.Unpack(r[0])
}

func (y YPod) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YPod) Restart() error {
	r := y.X.Call("Restart")
	return errors.Unpack(r[0])
}

func (y YPod) Scrub() {
	y.X.Call("Scrub")
}

func (y YPod) Signal(sig string) error {
	r := y.X.Call("Signal")
	return errors.Unpack(r[0])
}

func (y YPod) Start(opts *podman.PodStartOptions) error {
	r := y.X.Call("Start", opts)
	return errors.Unpack(r[0])
}

func (y YPod) Stop(opts *podman.PodStopOptions) error {
	r := y.X.Call("Stop", opts)
	return errors.Unpack(r[0])
}

func (y YPod) Unpause() error {
	r := y.X.Call("Unpause")
	return errors.Unpack(r[0])
}
