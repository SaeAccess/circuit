package makers

import (
	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/kit/pubsub"
)

var JoinType = reflect.TypeOf(YJoin{})

func init() {
	client.RegisterElementMaker(&joinElementMaker{
		client.NewBaseElementMaker("@join", reflect.TypeOf(YJoin{})),
	})
}

// implementation for a specific maker
type joinElementMaker struct {
	client.BaseElementMaker
}

func (b *joinElementMaker) Make(y anchor.YTerminal, arg any) (v any, err error) {
	v, err = y.Make(b.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if reflect.TypeOf(v) != b.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", b.Name())

	}

	// v can now be type asserted to t whithout error
	return YJoin{v.(pubsub.YSubscription)}, nil
}

func (c *joinElementMaker) Get(v any) any {
	return YJoin{v.(pubsub.YSubscription)}
}

func subscriptionStat(s pubsub.Stat) client.SubscriptionStat {
	return client.SubscriptionStat{
		Source:  s.Source,
		Pending: s.Pending,
		Closed:  s.Closed,
	}
}

type YJoin struct {
	pubsub.YSubscription
}

func (y YJoin) Peek() client.SubscriptionStat {
	return subscriptionStat(y.YSubscription.Peek())
}
