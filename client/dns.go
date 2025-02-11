// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

// NameserverStat encloses process state information.
type NameserverStat struct {

	// IP address of the nameserver
	Address string

	// Resource records resolved by this nameserver
	Records map[string][]string
}

type Nameserver interface {
	Set(rr string) error

	Unset(name string)

	// Peek asynchronously returns the current state of the server.
	Peek() NameserverStat

	PeekBytes() []byte

	// Scrub shuts down the nameserver and removes its circuit element.
	Scrub()
}
