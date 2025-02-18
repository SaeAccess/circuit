// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	// "fmt"

	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/tissue"
)

// An Anchor represents a location in the global anchor namespace of a circuit
// cluster. Anchors are named locations where the user can store and operate
// control primitives, called circuit elements. The anchor namespace hierarchy
// is represented in paths of the form
//
//	/X8817c114d4941522/hello/dolly
//
// The root anchor "/" represents the cluster abstractly and is the only
// anchor within which one cannot create elements or freely-named subanchors.
// The root anchor contains a dynamically changing set of sub-anchors that
// correspond to the live circuit servers in the cluster.
//
// Every anchor, other than "/", can be used to make, store and operate a
// circuit element (a process or a channel). Anchors are created on access, if
// not present, and are garbage-collected when not used or referenced.
// Therefore the interface allows users to access arbitrary paths without
// having to create them first.
type Anchor interface {

	// Addr returns the address of the circuit server hosting this anchor.
	Addr() string

	// ServerID returns the ID of the circuit server hosting this anchor.
	// The returned string will look like "X123..."
	ServerID() string

	// Walk traverses the anchor namespace, starting from this anchor along the path in walk.
	// Errors in communication or a missing circuit server condition are reported via panics.
	Walk(walk []string) Anchor

	// View returns the set of this anchor's sub-anchors.
	View() map[string]Anchor

	Make(typ reflect.Type, arg any) (any, error)

	// Get returns a handle for the circuit element (Chan, Proc, Subscription, Server, etc)
	// stored at this anchor, and nil otherwise.
	// Panics indicate that the server hosting the anchor and its element has already died.
	// Get() interface{}
	Get() any

	// Scrub aborts and abandons the circuit element stored at this anchor, if one is present.
	// If the hosting server is dead, a panic will be issued.
	Scrub()

	// Path returns the path to this anchor
	Path() string
}

// Alias these since they are both Subscription and would cause collision in ElementMaker map
type Leave Subscription
type Join Subscription

// Split breaks up an anchor path into components.
func Split(walk string) (r []string) {
	var j int
	for i, c := range walk {
		if c != '/' {
			continue
		}
		if i-j > 0 {
			r = append(r, walk[j:i])
		}
		j = i + 1
	}
	if len(walk)-j > 0 {
		r = append(r, walk[j:])
	}
	return
}

type terminal struct {
	y anchor.YTerminal
	k tissue.KinAvatar
}

func (t terminal) Addr() string {
	return t.k.X.Addr().String()
}

func (t terminal) ServerID() string {
	return t.k.ID.String()
}

func (t terminal) Walk(walk []string) Anchor {
	return terminal{y: t.y.Walk(walk), k: t.k}
}

func (t terminal) Path() string {
	return t.y.Path()
}

func (t terminal) View() map[string]Anchor {
	v := t.y.View()
	w := make(map[string]Anchor)
	for name, y := range v {
		w[name] = terminal{y: y, k: t.k}
	}
	return w
}

// Single Make method for anchor impl. The specified type argument should be a registered
// interface that the terminal knows how to make
func (t terminal) Make(typ reflect.Type, arg any) (any, error) {
	// type should be interface type that is registered as a maker type
	if typ.Kind() != reflect.Interface {
		return nil, fmt.Errorf("Make requires an interface reflected type that was registered as a maker")
	}

	// Get the maker for the specified type
	maker, ok := GetElementMaker(typ)
	if !ok {
		return nil, fmt.Errorf("no element maker for type=%v", typ)
	}

	el, err := maker.Make(t.y, arg)
	if err != nil {
		return nil, err
	}

	// Validate the type, maybe maker didnt
	if !reflect.TypeOf(el).Implements(typ) {
		return nil, fmt.Errorf("mismatch element maker types, %s != %s", reflect.TypeOf(el).Name(), typ.Name())
	}

	return el, nil
}

func (t terminal) Get() any {
	_, y := t.y.Get()
	if y == nil {
		return nil
	}

	return y
}

func (t terminal) Scrub() {
	t.y.Scrub()
}

func GetElement[T any](a Anchor, path []string) T {
	el := a.Walk(path).Get()
	v, ok := el.(T)
	if !ok {
		var t T
		return t
	}

	return v
}

func MakeElement[T any](a Anchor, arg any, path ...string) (T, error) {
	mt := reflect.TypeOf((*T)(nil)).Elem()
	el, err := a.Walk(path).Make(mt, arg)
	var t T
	if err != nil {
		return t, err
	}

	v, ok := el.(T)
	if !ok {
		return t, err
	}

	return v, nil
}
