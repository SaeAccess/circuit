// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"fmt"

	"github.com/gocircuit/circuit/use/circuit"
	xerrors "github.com/gocircuit/circuit/use/errors"
)

func init() {
	circuit.RegisterValue(XTerminal{})
}

// XTerminal represents the server side stub that receives the remote call from
// YTerminal
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

// YTerminal represents client side stub used for making remote calls.
type YTerminal struct {
	X circuit.X
}

func (y YTerminal) Path() string {
	return y.X.Call("Path")[0].(string)
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

// Make call the server to create the anchor and element
func (y YTerminal) Make(kind string, arg interface{}) (yelm interface{}, err error) {
	r := y.X.Call("Make", kind, arg)
	if err = xerrors.Unpack(r[1]); err != nil {
		return nil, err
	}

	x, ok := r[0].(circuit.X)
	if !ok {
		return nil, fmt.Errorf("element not created, kind=%s, got %v", kind, r[0])
	}

	factory, ok := efRepo.GetYF(kind)
	if !ok {
		return nil, fmt.Errorf("element kind=%s not known", kind)
	}

	return factory(x)
}

// Get calls the server to get the remote anchor element
func (y YTerminal) Get() (kind string, yelm interface{}) {
	r := y.X.Call("Get")
	kind = r[0].(string)
	x, ok := r[1].(circuit.X)
	if !ok {
		return "", nil
	}

	factory, ok := efRepo.GetYF(kind)
	if !ok {
		return "", fmt.Errorf("element kind=%s not known", kind)
	}

	v, _ := factory(x)
	return kind, v
}

func (y YTerminal) Scrub() {
	y.X.Call("Scrub")
}
