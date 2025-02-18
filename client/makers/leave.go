package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
)

func init() {
	client.RegisterElementMaker(&leaveElementMaker{
		client.NewBaseElementMaker("@leave", LeaveType),
	})
}

var LeaveType = reflect.TypeOf((*client.Leave)(nil)).Elem()

// implementation for a specific maker
type leaveElementMaker struct {
	client.BaseElementMaker
}

func (b *leaveElementMaker) Make(y anchor.YTerminal, arg any) (v any, err error) {
	v, err = y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if !reflect.TypeOf(v).Implements(b.Type()) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), LeaveType)
	}

	// v can now be type asserted to t whithout error
	return v, nil //YLeave{v.(pubsub.YSubscription)}, nil
}

// type YLeave struct {
// 	pubsub.YSubscription
// }

// func (y YLeave) Peek() client.SubscriptionStat {
// 	return subscriptionStat(y.YSubscription.Peek())
// }
