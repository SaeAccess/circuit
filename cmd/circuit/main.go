// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"log"
	"os"
)

// TODO fix Makefile so that main can be moved out and package renamed to circuit
func main() {
	if err := Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
