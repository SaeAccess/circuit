package wasm

import (
	"io"

	ws "github.com/gocircuit/circuit/client/wasm"
	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XWasm{})
}

// XWasm repreents the skeleton for the server side Wasm container
type XWasm struct {
	Wasm
}

func (x XWasm) Wait() (*ws.Status, error) {
	stat, err := x.Wasm.Wait()
	return stat, errors.Pack(err)
}

func (x XWasm) Signal(sig string) error {
	return errors.Pack(x.Wasm.Signal(sig))
}

func (x XWasm) Stdin() circuit.X {
	return xio.NewXWriteCloser(x.Wasm.Stdin())
}

func (x XWasm) Stdout() circuit.X {
	return xio.NewXReadCloser(x.Wasm.Stdout())
}

func (x XWasm) Stderr() circuit.X {
	return xio.NewXReadCloser(x.Wasm.Stderr())
}

func (x XWasm) Peek() (*ws.Status, error) {
	stat, err := x.Wasm.Peek()
	return stat, errors.Pack(err)
}

func (x XWasm) PeekBytes() []byte {
	return x.Wasm.PeekBytes()
}

type YWasm struct {
	X circuit.X
}

func (y YWasm) Wait() (stat *ws.Status, err error) {
	r := y.X.Call("Wait")
	stat, _ = r[0].(*ws.Status)
	return stat, errors.Unpack(r[1])
}

func (y YWasm) Signal(sig string) error {
	r := y.X.Call("Signal", sig)
	return errors.Unpack(r[0])
}

func (y YWasm) Scrub() {
	y.X.Call("Scrub")
}

func (y YWasm) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YWasm) Peek() (stat *ws.Status, err error) {
	r := y.X.Call("Peek")
	stat, _ = r[0].(*ws.Status)
	return stat, errors.Unpack(r[1])
}

func (y YWasm) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YWasm) Stdin() io.WriteCloser {
	return xio.NewYWriteCloser(y.X.Call("Stdin")[0])
}

func (y YWasm) Stdout() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stdout")[0])
}

func (y YWasm) Stderr() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stderr")[0])
}
