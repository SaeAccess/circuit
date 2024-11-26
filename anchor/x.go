// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"errors"

	"github.com/gocircuit/circuit/element/dns"
	"github.com/gocircuit/circuit/element/docker"
	"github.com/gocircuit/circuit/element/proc"
	srv "github.com/gocircuit/circuit/element/server"
	"github.com/gocircuit/circuit/element/valve"
	"github.com/gocircuit/circuit/kit/pubsub"
	"github.com/gocircuit/circuit/use/circuit"
	xerrors "github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XTerminal{})
}

type XTerminal struct {
	t *Terminal
}

func (x XTerminal) Path() string {
	return x.t.Path()
}

func (x XTerminal) Walk(walk []string) circuit.X {
	t := x.t.Walk(walk)
	if t == nil {
		return nil
	}
	return circuit.Ref(XTerminal{t})
}

func (x XTerminal) View() map[string]circuit.X {
	u := make(map[string]circuit.X)
	for p, q := range x.t.View() {
		u[p] = circuit.Ref(XTerminal{q})
	}
	return u
}

func (x XTerminal) Make(kind string, arg interface{}) (xelm circuit.X, err error) {
	elm, err := x.t.Make(kind, arg)
	if err != nil {
		return nil, xerrors.Pack(err)
	}
	return elm.X(), nil
}

func (x XTerminal) Get() (string, circuit.X) {
	kind, elm := x.t.Get()
	if elm == nil {
		return "", nil
	}
	return kind, elm.X()
}

func (x XTerminal) Scrub() {
	x.t.Scrub()
}

// YTerminal…
type YTerminal struct {
	X circuit.X
}

func (y YTerminal) Walk(walk []string) YTerminal {
	return YTerminal{
		y.X.Call("Walk", walk)[0].(circuit.X),
	}
}

func (y YTerminal) View() map[string]YTerminal {
	u := make(map[string]YTerminal)
	r := y.X.Call("View")
	for n, x := range r[0].(map[string]circuit.X) {
		u[n] = YTerminal{x}
	}
	return u
}

func (y YTerminal) Make(kind string, arg interface{}) (yelm interface{}, err error) {
	r := y.X.Call("Make", kind, arg)
	if err = xerrors.Unpack(r[1]); err != nil {
		return nil, err
	}
	switch kind {
	case Chan:
		return valve.YValve{X: r[0].(circuit.X)}, nil
	case Proc:
		return proc.YProc{X: r[0].(circuit.X)}, nil
	case Docker:
		return docker.YContainer{X: r[0].(circuit.X)}, nil
	case Nameserver:
		return dns.YNameserver{X: r[0].(circuit.X)}, nil
	case OnJoin:
		return pubsub.YSubscription{X: r[0].(circuit.X)}, nil
	case OnLeave:
		return pubsub.YSubscription{X: r[0].(circuit.X)}, nil
	}
	return nil, errors.New("element kind not supported")
}

func (y YTerminal) Get() (kind string, yelm interface{}) {
	r := y.X.Call("Get")
	kind = r[0].(string)
	switch kind {
	case Server:
		return Server, srv.YServer{X: r[1].(circuit.X)}
	case Chan:
		return Chan, valve.YValve{X: r[1].(circuit.X)}
	case Proc:
		return Proc, proc.YProc{X: r[1].(circuit.X)}
	case Nameserver:
		return Nameserver, dns.YNameserver{X: r[1].(circuit.X)}
	case Docker:
		return Docker, docker.YContainer{X: r[1].(circuit.X)}
	case OnJoin:
		return OnJoin, pubsub.YSubscription{X: r[1].(circuit.X)}
	case OnLeave:
		return OnLeave, pubsub.YSubscription{X: r[1].(circuit.X)}
	}
	return "", nil
}

func (y YTerminal) Scrub() {
	y.X.Call("Scrub")
}

func (y YTerminal) Path() string {
	return y.X.Call("Path")[0].(string)
}
