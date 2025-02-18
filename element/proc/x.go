// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package proc

import (
	"io"

	"github.com/gocircuit/circuit/client"
	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XProc{})
}

type XProc struct {
	Proc
}

// func unpack(stat Stat) Stat {
// 	stat.Exit = errors.Unpack(stat.Exit)
// 	return stat
// }

// func pack(stat Stat) Stat {
// 	stat.Exit = errors.Pack(stat.Exit)
// 	return stat
// }

func (x XProc) Wait() (client.ProcStat, error) {
	stat, err := x.Proc.Wait()
	stat.Exit = errors.Pack(stat.Exit)
	return stat, errors.Pack(err)
}

func (x XProc) Signal(sig string) error {
	return errors.Pack(x.Proc.Signal(sig))
}

func (x XProc) Stdin() circuit.X {
	return xio.NewXWriteCloser(x.Proc.Stdin())
}

func (x XProc) Stdout() circuit.X {
	return xio.NewXReadCloser(x.Proc.Stdout())
}

func (x XProc) Stderr() circuit.X {
	return xio.NewXReadCloser(x.Proc.Stderr())
}

func (x XProc) Peek() client.ProcStat {
	ps := x.Proc.Peek()
	ps.Exit = errors.Pack(ps.Exit)
	return ps
}

func (x XProc) PeekBytes() []byte {
	return x.Proc.PeekBytes()
}

type YProc struct {
	X circuit.X
}

func (y YProc) Wait() (client.ProcStat, error) {
	r := y.X.Call("Wait")
	ps := r[0].(client.ProcStat)
	ps.Exit = errors.Unpack(ps.Exit)
	return ps, errors.Unpack(r[1])
}

func (y YProc) Signal(sig string) error {
	r := y.X.Call("Signal", sig)
	return errors.Unpack(r[0])
}

func (y YProc) Scrub() {
	y.X.Call("Scrub")
}

func (y YProc) GetEnv() []string {
	return y.X.Call("GetEnv")[0].([]string)
}

func (y YProc) GetCmd() client.Cmd {
	return y.X.Call("GetCmd")[0].(client.Cmd)
}

func (y YProc) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YProc) Peek() client.ProcStat {
	ps := y.X.Call("Peek")[0].(client.ProcStat)
	ps.Exit = errors.Unpack(ps.Exit)
	return ps
}

func (y YProc) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YProc) Stdin() io.WriteCloser {
	return xio.NewYWriteCloser(y.X.Call("Stdin")[0])
}

func (y YProc) Stdout() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stdout")[0])
}

func (y YProc) Stderr() io.ReadCloser {
	return xio.NewYReadCloser(y.X.Call("Stderr")[0])
}
