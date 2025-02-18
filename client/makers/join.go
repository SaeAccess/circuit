package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
)

var JoinType = reflect.TypeOf((*client.Join)(nil)).Elem()

func init() {
	client.RegisterElementMaker(&joinElementMaker{
		client.NewBaseElementMaker("@join", JoinType),
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
	if !reflect.TypeOf(v).Implements(JoinType) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), JoinType)
	}

	// v can now be type asserted to t whithout error
	return v, nil //YJoin{v.(pubsub.YSubscription)}, nil
}

// func subscriptionStat(s pubsub.Stat) client.SubscriptionStat {
// 	return client.SubscriptionStat{
// 		Source:  s.Source,
// 		Pending: s.Pending,
// 		Closed:  s.Closed,
// 	}
// }

// type YJoin struct {
// 	pubsub.YSubscription
// }

// func (y YJoin) Peek() client.SubscriptionStat {
// 	return subscriptionStat(y.YSubscription.Peek())
// }
