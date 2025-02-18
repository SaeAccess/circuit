package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
)

func init() {
	client.RegisterElementMaker(&dnsElementMaker{
		client.NewBaseElementMaker("dns", NameserverType),
	})
}

var NameserverType = reflect.TypeOf((*client.Nameserver)(nil)).Elem()

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
	if !reflect.TypeOf(v).Implements(NameserverType) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), NameserverType)
	}

	// v can now be type asserted to t whithout error
	return v, nil //yNameserver{v.(dns.YNameserver)}, nil
}

// func nameserverStat(s dns.Stat) client.NameserverStat {
// 	return client.NameserverStat{
// 		Address: s.Address,
// 		Records: s.Records,
// 	}
// }

// type yNameserver struct {
// 	dns.YNameserver
// }

// func (y yNameserver) Peek() client.NameserverStat {
// 	return nameserverStat(y.YNameserver.Peek())
// }
