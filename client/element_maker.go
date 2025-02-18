package client

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gocircuit/circuit/anchor"
	"github.com/pkg/errors"
)

// ElementMaker allows for extending the client interface with new anchor elements
// without the need to modify the client api.  This is along with the element maker
// repository supports registration of new elment maker types.
type ElementMaker interface {
	Name() string
	Type() reflect.Type
	Make(y anchor.YTerminal, arg any) (any, error) // TODO maybe could take many args
	Is(t reflect.Type) bool
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

func (e *MakerRepo) Find(t reflect.Type) (m ElementMaker, b bool) {
	e.sm.Range(func(key, value any) bool {
		if value.(ElementMaker).Is(t) {
			m = value.(ElementMaker)
			b = true
			return false
		}

		return true
	})

	return
}

var (
	mr              = MakerRepo{}
	ErrMismatchType = fmt.Errorf("client/circuit type mismatch")
)

func RegisterElementMaker(m ElementMaker) {
	//log.Printf("Registering element maker: name=%s type=%v", m.Name(), m.Type())
	mr.Add(m)
}

func GetElementMaker(t reflect.Type) (ElementMaker, bool) {
	return mr.Get(t)
}

func NewBaseElementMaker(name string, t reflect.Type) BaseElementMaker {
	return BaseElementMaker{name: name, t: t}
}

func FindElementMaker(v any) ElementMaker {
	if m, ok := mr.Find(reflect.TypeOf(v)); ok {
		return m
	}

	return nil
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
	if b.Type().Kind() == reflect.Interface {
		if !reflect.TypeOf(v).Implements(b.Type()) {
			return nil, errors.Wrapf(ErrMismatchType, "kind=%v", b.Name())
		}
	} else {
		if reflect.TypeOf(v) != b.Type() {
			return nil, errors.Wrapf(ErrMismatchType, "kind=%v", b.Name())
		}

	}

	return v, nil // v can now be type asserted to t whithout error
}
func (b *BaseElementMaker) Is(t reflect.Type) bool {
	return t.Implements(b.t)
}
