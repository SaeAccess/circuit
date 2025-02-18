// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package server

import (
	"io"
	"time"
)

// Server…
// All methods panic if the hosting circuit server dies.
type Server interface {
	Profile(string) (io.ReadCloser, error)
	Peek() ServerStat
	PeekBytes() []byte
	Rejoin(string) error
	Suicide()
}

// ServerStat encloses subscription state information.
type ServerStat struct {
	Addr   string    `json:"addr,omitempty"`
	Joined time.Time `json:"joined,omitempty"`

	// TODO - put some runtime metrics for the server here and also add a  map for adding metadata
	// for the server.  This can be used for scheduling elements and runtime behavior reporting.

	// TODO access policies, how to implement
}
