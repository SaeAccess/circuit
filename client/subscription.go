// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

// SubscriptionStat encloses subscription state information.
type SubscriptionStat struct {

	// Name of event source
	Source string

	// Pending equals the number of messages waiting to be consumed.
	Pending int

	// Closed is true if the publisher stream has marked an end.
	Closed bool
}

// Subscription provides access to a circuit subscription element.
// All methods panic if the hosting circuit server dies.
type Subscription interface {

	// Consume blocks until the next message is available on the channel.
	Consume() (interface{}, bool)

	// Peek asynchronously returns the current state of the process.
	Peek() SubscriptionStat

	// Scrub abandons the circuit process element, without affecting the underlying OS process.
	Scrub()
}
