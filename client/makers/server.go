package makers

import (
	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/server"
)

func init() {
	client.RegisterElementMaker(&serverElementMaker{
		client.NewBaseElementMaker("server", reflect.TypeOf(server.YServer{})),
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
	if reflect.TypeOf(v) != b.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", b.Name())
	}

	// v can now be type asserted to t whithout error
	return ysrvSrv{v.(server.YServer)}, nil
}

func (c *serverElementMaker) Get(v any) any {
	return ysrvSrv{v.(server.YServer)}
}

func srvStat(s server.Stat) client.ServerStat {
	return client.ServerStat{
		Addr:   s.Addr,
		Joined: s.Joined,
	}
}

type ysrvSrv struct {
	server.YServer
}

func (y ysrvSrv) Peek() client.ServerStat {
	return srvStat(y.YServer.Peek())
}
