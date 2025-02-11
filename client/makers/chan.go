package makers

import (
	"errors"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/valve"
)

var ChanType = reflect.TypeOf(valve.YValve{})

func init() {
	client.RegisterElementMaker(&chanElementMaker{
		client.NewBaseElementMaker("chan", reflect.TypeOf(valve.YValve{})),
	})
}

// implementation for a specific maker
type chanElementMaker struct {
	client.BaseElementMaker
}

func (b *chanElementMaker) Make(y anchor.YTerminal, arg any) (any, error) {
	v, err := y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if reflect.TypeOf(v) != b.Type() {
		return nil, errors.New("mismatch chan")
	}

	// v can now be type asserted to t whithout error
	return yvalveChan{v.(valve.YValve)}, nil
}

func (c *chanElementMaker) Get(v any) any {
	return yvalveChan{v.(valve.YValve)}
}

func retypeChanStat(s valve.Stat) client.ChanStat {
	return client.ChanStat{
		Cap:     s.Cap,
		Closed:  s.Closed,
		Aborted: s.Aborted,
		NumSend: s.NumSend,
		NumRecv: s.NumRecv,
	}
}

type yvalveChan struct {
	valve.YValve
}

func (y yvalveChan) Stat() client.ChanStat {
	return retypeChanStat(y.YValve.Stat())
}
