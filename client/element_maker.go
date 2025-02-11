package client

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gocircuit/circuit/anchor"
)

// ElementMaker allows for extending the client interface with new anchor elements
// without the need to modify the client api.  This is along with the element maker
// repository supports registration of new elment maker types.
type ElementMaker interface {
	Name() string
	Type() reflect.Type
	Make(y anchor.YTerminal, arg any) (any, error) // TODO maybe could take many args
	Get(any) any
}

// MakerRepo maintains the registered element makers
type MakerRepo struct {
	sm sync.Map // map[reflect.Type]ElementMaker
}

func (e *MakerRepo) Add(n ElementMaker) {
	e.sm.Store(n.Type(), n)
}

func (e *MakerRepo) Get(t reflect.Type) (ElementMaker, bool) {
	em, ok := e.sm.Load(t)
	if !ok {
		return nil, ok
	}

	return em.(ElementMaker), true
}

var (
	mr = MakerRepo{}
)

func RegisterElementMaker(m ElementMaker) {
	mr.Add(m)
}

func GetElementMaker(t reflect.Type) (ElementMaker, bool) {
	return mr.Get(t)
}

func NewBaseElementMaker(name string, t reflect.Type) BaseElementMaker {
	return BaseElementMaker{name: name, t: t}
}

type BaseElementMaker struct {
	name string
	t    reflect.Type
}

func (b BaseElementMaker) Name() string       { return b.name }
func (b BaseElementMaker) Type() reflect.Type { return b.t }
func (b *BaseElementMaker) Make(y anchor.YTerminal, arg any) (any, error) {
	v, err := y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if reflect.TypeOf(v) != b.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", b.Name())
	}

	return v, nil // v can now be type asserted to t whithout error
}

// ===============

// func use() {
// 	nw := a.Make(reflect.TypeOf(network{}), networkArg{}).(network)
// }
