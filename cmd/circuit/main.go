// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:S
//   2013 Petar Maymounkov <p@gocircuit.org>

// This package provides the executable program for the resource-sharing circuit app
package main

import (
	"log"
	"os"

	"github.com/gocircuit/circuit/cmd"
	_ "github.com/gocircuit/circuit/element/dns"
	_ "github.com/gocircuit/circuit/element/docker"
	_ "github.com/gocircuit/circuit/element/podman/container"
	_ "github.com/gocircuit/circuit/element/podman/network"
	_ "github.com/gocircuit/circuit/element/podman/pod"
	_ "github.com/gocircuit/circuit/element/podman/volume"
	_ "github.com/gocircuit/circuit/element/proc"
	_ "github.com/gocircuit/circuit/element/server"
	_ "github.com/gocircuit/circuit/element/valve"
	_ "github.com/gocircuit/circuit/element/wasm"
)

// TODO fix Makefile so that main can be moved out and package renamed to circuit
func main() {
	if err := cmd.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
