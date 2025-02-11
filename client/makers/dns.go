package makers

import (
	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/dns"
)

func init() {
	client.RegisterElementMaker(&dnsElementMaker{
		client.NewBaseElementMaker("dns", reflect.TypeOf(dns.YNameserver{})),
	})
}

var NameserverType = reflect.TypeOf(dns.YNameserver{})

// implementation for a specific maker
type dnsElementMaker struct {
	client.BaseElementMaker
}

func (b *dnsElementMaker) Make(y anchor.YTerminal, arg any) (any, error) {
	v, err := y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if reflect.TypeOf(v) != b.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", b.Name())

	}

	// v can now be type asserted to t whithout error
	return yNameserver{v.(dns.YNameserver)}, nil
}

func (c *dnsElementMaker) Get(v any) any {
	return yNameserver{v.(dns.YNameserver)}
}

func nameserverStat(s dns.Stat) client.NameserverStat {
	return client.NameserverStat{
		Address: s.Address,
		Records: s.Records,
	}
}

type yNameserver struct {
	dns.YNameserver
}

func (y yNameserver) Peek() client.NameserverStat {
	return nameserverStat(y.YNameserver.Peek())
}
