package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
)

var ChanType = reflect.TypeOf((*client.Chan)(nil)).Elem() //valve.YValve{})

func init() {
	client.RegisterElementMaker(&chanElementMaker{
		client.NewBaseElementMaker("chan", ChanType),
	})
}

// implementation for a specific maker
type chanElementMaker struct {
	client.BaseElementMaker
}

func (b *chanElementMaker) Make(y anchor.YTerminal, arg any) (v any, err error) {
	v, err = y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	if !reflect.TypeOf(v).Implements(ChanType) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), ChanType)
	}

	// v can now be type asserted to t whithout error
	return v, nil // yvalveChan{v.(valve.YValve)}, nil
}

// func (c *chanElementMaker) Get(v any) any {
// 	return yvalveChan{v.(valve.YValve)}
// }

// func retypeChanStat(s valve.Stat) client.ChanStat {
// 	return client.ChanStat{
// 		// Cap:     s.Cap,
// 		Closed:  s.Closed,
// 		Aborted: s.Aborted,
// 		NumSend: s.NumSend,
// 		NumRecv: s.NumRecv,
// 	}
// }

// type yvalveChan struct {
// 	valve.YValve
// }

// func (y yvalveChan) Stat() client.ChanStat {
// 	return y.YValve.Stat() //retypeChanStat(y.YValve.Stat())
// }
