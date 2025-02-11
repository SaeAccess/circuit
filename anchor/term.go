// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package anchor

import (
	"errors"
	"fmt"
	"log"
	"sync"

	srv "github.com/gocircuit/circuit/element/server"
	"github.com/gocircuit/circuit/kit/pubsub"
	"github.com/gocircuit/circuit/use/circuit"
)

// Element are the values placed in anchor nodes
type Element interface {
	Scrub()
	X() circuit.X
}

// Peeker provides a generic mechanism for peeking at any anchor type.  Each anchor
// must implement the Peeker interface that converts its actual peek state to a byte array
type Peeker interface {
	PeekBytes() []byte
}

type ElementFactory func(t *Terminal, arg any) (Element, error)
type YFactory func(x circuit.X) (any, error)

const (
	Server     = "server"
	Chan       = "chan"
	Proc       = "proc"
	Docker     = "docker"
	Nameserver = "dns"
	OnJoin     = "@join"
	OnLeave    = "@leave"

	// wasm
	Wasm = "wasm"

	Pod = "pod"

	// oci container support by podman/containerd
	Container = "container"

	// Network
	Network = "network"

	// Volume
	Volume = "volume"
)

func init() {
	registerServer()
	registerOnJoin()
	registerOnLeave()
}

var efRepo = newElementFactoryRepo()

// RegisterElement registers an element factory
func RegisterElement(kind string, efactory ElementFactory, yfactory YFactory) {
	efRepo.Register(kind, efactory, yfactory)
}

// Terminal presents a facade to *Anchor with added element manipulation methods
type Terminal struct {
	genus  Genus
	anchor *Anchor
}

type Genus interface {
	NewArrivals() pubsub.Consumer
	NewDepartures() pubsub.Consumer
}

// NewTerm create the root node of a new anchor file system.
func NewTerm(name string, genus Genus) (*Terminal, circuit.PermX) {
	t := &Terminal{
		genus:  genus,
		anchor: newAnchor(nil, name).use(),
	}
	return t, circuit.PermRef(XTerminal{t})
}

func (t *Terminal) carrier() *Anchor {
	return t.anchor
}

func (t *Terminal) Walk(walk []string) *Terminal {
	return &Terminal{
		genus:  t.genus,
		anchor: t.carrier().Walk(walk),
	}
}

func (t *Terminal) Path() string {
	return t.carrier().Path()
}

func (t *Terminal) View() map[string]*Terminal {
	r := make(map[string]*Terminal)
	for n, a := range t.carrier().View() {
		r[n] = &Terminal{
			genus:  t.genus,
			anchor: a,
		}
	}
	return r
}

type urn struct {
	kind string
	elem Element // valve.Valve, proc.Proc, etc
}

func (t *Terminal) Attach(kind string, elm Element) {
	if kind != Server {
		panic(0)
	}
	log.Printf("Attaching %s as %s", t.carrier().Path(), kind)
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	if t.carrier().Get() != nil {
		panic(0)
	}
	u := &urn{
		kind: kind,
		elem: elm,
	}
	t.carrier().Set(u)
}

func (t *Terminal) Make(kind string, arg interface{}) (eleme Element, err error) {
	log.Printf("Detaching %s", t.carrier().Path())
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	if t.carrier().Get() != nil {
		return nil, errors.New("anchor already has an element")
	}

	// get the element factory
	factory, ok := efRepo.GetEF(kind)
	if !ok {
		return nil, fmt.Errorf("element kind not known, kind=%s", kind)
	}

	// create the element
	eleme, err = factory(t, arg)
	if err != nil {
		return nil, err
	}

	u := &urn{
		kind: kind,
		elem: eleme,
	}
	t.carrier().Set(u)
	return u.elem, nil
}

func (t *Terminal) Get() (string, Element) {
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	v := t.carrier().Get()
	if v == nil {
		return "", nil
	}
	return v.(*urn).kind, v.(*urn).elem
}

func (t *Terminal) Scrub() {
	log.Printf("Scrubbing %s", t.carrier().Path())
	t.carrier().TxLock()
	defer t.carrier().TxUnlock()
	u, ok := t.carrier().Get().(*urn)
	if !ok {
		return
	}
	if _, ok := u.elem.(srv.Server); ok {
		return // Cannot scrub server anchors
	}
	u.elem.Scrub()
	t.carrier().Set(nil)
}

type elementFactoryRepo struct {
	ef sync.Map // map[string]ElementFactory
}

type factoryItem struct {
	ef ElementFactory
	yf YFactory
}

func newElementFactoryRepo() *elementFactoryRepo {
	return &elementFactoryRepo{}
}

func (r *elementFactoryRepo) Register(kind string, efactory ElementFactory, yfactory YFactory) {
	// check if already registered
	if _, ok := r.GetEF(kind); ok {
		panic(fmt.Sprintf("%s element already registered", kind))
	}
	r.ef.Store(kind, factoryItem{efactory, yfactory})
}

func (r *elementFactoryRepo) GetEF(kind string) (ElementFactory, bool) {
	v, ok := r.ef.Load(kind)
	if !ok {
		return nil, false
	}
	return v.(factoryItem).ef, true
}

func (r *elementFactoryRepo) GetYF(kind string) (YFactory, bool) {
	v, ok := r.ef.Load(kind)
	if !ok {
		return nil, false
	}
	return v.(factoryItem).yf, true
}

func registerServer() {
	RegisterElement(Server,
		func(t *Terminal, arg any) (Element, error) {
			panic("element factory not implemented for server")
		},

		func(x circuit.X) (any, error) {
			return srv.YServer{X: x}, nil
		})
}

func registerOnJoin() {
	RegisterElement(OnJoin,
		func(t *Terminal, arg any) (Element, error) {
			return t.genus.NewArrivals(), nil
		},

		func(x circuit.X) (any, error) {
			return pubsub.YSubscription{X: x}, nil
		})
}

func registerOnLeave() {
	RegisterElement(OnLeave,
		func(t *Terminal, arg any) (Element, error) {
			return t.genus.NewDepartures(), nil
		},

		func(x circuit.X) (any, error) {
			return pubsub.YSubscription{X: x}, nil
		})
}
