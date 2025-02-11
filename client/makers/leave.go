package makers

import (
	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/kit/pubsub"
)

func init() {
	client.RegisterElementMaker(&leaveElementMaker{
		client.NewBaseElementMaker("@leave", reflect.TypeOf(YLeave{})),
	})
}

var LeaveType = reflect.TypeOf(YLeave{})

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
	if reflect.TypeOf(v) != b.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", b.Name())
	}

	// v can now be type asserted to t whithout error
	return YLeave{v.(pubsub.YSubscription)}, nil
}

func (c *leaveElementMaker) Get(v any) any {
	return YLeave{v.(pubsub.YSubscription)}
}

type YLeave struct {
	pubsub.YSubscription
}

func (y YLeave) Peek() client.SubscriptionStat {
	return subscriptionStat(y.YSubscription.Peek())
}
