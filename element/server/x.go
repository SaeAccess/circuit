// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package server

import (
	// "fmt"
	"io"

	cli "github.com/gocircuit/circuit/client/server"
	xio "github.com/gocircuit/circuit/kit/x/io"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XServer{})
}

// XServer…
type XServer struct {
	*server
}

func (x XServer) Profile(name string) (circuit.X, error) {
	r, err := x.server.Profile(name)
	if err != nil {
		return nil, errors.Pack(err)
	}
	return xio.NewXReadCloser(r), nil
}

func (x XServer) Rejoin(addr string) error {
	return errors.Pack(x.server.Rejoin(addr))
}

// YServer…
type YServer struct {
	X circuit.X
}

func (y YServer) Profile(name string) (rc io.ReadCloser, err error) {
	r := y.X.Call("Profile", name)
	if err := errors.Unpack(r[1]); err != nil {
		return nil, err
	}
	return xio.NewYReadCloser(r[0]), nil
}

func (y YServer) Peek() cli.ServerStat {
	return y.X.Call("Peek")[0].(cli.ServerStat)
}

func (y YServer) PeekBytes() []byte {
	return y.X.Call("PeekBytes")[0].([]byte)
}

func (y YServer) Rejoin(addr string) error {
	return errors.Unpack(y.X.Call("Rejoin", addr)[0])
}

func (y YServer) IsDone() bool {
	return y.X.Call("IsDone")[0].(bool)
}

func (y YServer) Scrub() {
	y.X.Call("Scrub")
}

func (y YServer) Suicide() {
	y.X.Call("Suicide")
}
