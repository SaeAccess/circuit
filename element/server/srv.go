// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"runtime/pprof"
	"time"

	"github.com/gocircuit/circuit/anchor"
	cli "github.com/gocircuit/circuit/client/server"
	"github.com/gocircuit/circuit/kit/interruptible"
	"github.com/gocircuit/circuit/tissue"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/gocircuit/circuit/use/n"
)

type Server interface {
	cli.Server
	Scrub()
	X() circuit.X
}

// server
type server struct {
	addr   string
	kin    *tissue.Kin
	joined time.Time
}

func init() {
	anchor.RegisterElement("server", ef, yf)
}

func New(kin *tissue.Kin) Server {
	return &server{
		addr:   kin.Avatar().X.Addr().String(),
		kin:    kin,
		joined: time.Now(),
	}
}

func (s *server) Rejoin(addr string) error {
	a, err := n.ParseAddr(addr)
	if err != nil {
		return err
	}
	return s.kin.ReJoin(a)
}

func (s *server) Suicide() {
	os.Exit(0)
}

func (s *server) Profile(name string) (io.ReadCloser, error) {
	p := pprof.Lookup(name)
	if p == nil {
		return nil, errors.New("no profile")
	}
	r, w := interruptible.Pipe()
	go func() {
		p.WriteTo(w, 1)
		w.Write([]byte("•••\n"))
		w.Close()
	}()
	return r, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

func (s *server) Peek() cli.ServerStat {
	return cli.ServerStat{
		Addr:   s.addr,
		Joined: s.joined,
	}
}

func (s *server) PeekBytes() []byte {
	b, _ := json.MarshalIndent(s.Peek(), "", "\t")
	return b
}

func (s *server) IsDone() bool {
	return false
}

func (s *server) Scrub() {}

func (s *server) X() circuit.X {
	return circuit.Ref(XServer{s})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	panic("element factory not implemented for server")
}

func yf(x circuit.X) (any, error) {
	return YServer{X: x}, nil
}
