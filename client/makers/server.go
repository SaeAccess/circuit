package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
)

var ServerType = reflect.TypeOf((*client.Server)(nil)).Elem()

func init() {
	client.RegisterElementMaker(&serverElementMaker{
		client.NewBaseElementMaker("server", ServerType),
	})
}

// implementation for a specific maker
type serverElementMaker struct {
	client.BaseElementMaker
}

func (b *serverElementMaker) Make(y anchor.YTerminal, arg any) (v any, err error) {
	v, err = y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if !reflect.TypeOf(v).Implements(b.Type()) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), ServerType)
	}

	// v can now be type asserted to t whithout error
	return v, nil //ysrvSrv{v.(server.YServer)}, nil
}

// func srvStat(s server.Stat) client.ServerStat {
// 	return client.ServerStat{
// 		Addr:   s.Addr,
// 		Joined: s.Joined,
// 	}
// }

// type ysrvSrv struct {
// 	server.YServer
// }

// func (y ysrvSrv) Peek() client.ServerStat {
// 	return srvStat(y.YServer.Peek())
// }
