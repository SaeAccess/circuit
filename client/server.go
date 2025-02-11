// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"io"
	"time"
)

// ServerStat encloses subscription state information.
type ServerStat struct {
	Addr   string
	Joined time.Time

	// TODO - put some runtime metrics for the server here and also add a  map for adding metadata
	// for the server.  This can be used for scheduling elements and runtime behavior reporting.

}

// Server…
// All methods panic if the hosting circuit server dies.
type Server interface {
	Profile(string) (io.ReadCloser, error)
	Peek() ServerStat
	Rejoin(string) error
	Suicide()
}
