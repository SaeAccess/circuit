package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gocircuit/circuit/client"
	_ "github.com/gocircuit/circuit/client/makers"
	"github.com/gocircuit/circuit/client/podman"
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

// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

/*
Virus is a simple, mildly-resilient to failure mechanism that runs around a cluster
and delivers its payload (a container execution), in a self-sustained fashion.

The virus mechanism contains two parts: a payload container and a nucleus process.

The payload can be any container, available via a container registery

The nucleus is embodied in this circuit application. It executes the payload container
on a random machine in the circuit cluster. Then it installs itself on a machine different
from that of the payload container, and proceeds to watch the payload until it dies.

When the payload dies, the nucleus executes a new payload instance on another
randomly chosen host, and replaces itself with a new nucleus process on yet
another new random host. And so on ...
*/

// pickServer returns the root anchor of a randomly-chosen circuit server in the cluster.
func pickServer(c *client.Client) client.Anchor {
	for _, r := range c.View() {
		return r
	}
	panic(0)
}

// waitForPayloadDeath blocks until the payload container stored at anchor exits, for whatever reason.
// anchor is the anchor path used by the virus logic.
func waitFotPayloadDeath(c *client.Client, myAnchor, payloadAnchor string, epoch int) (recov interface{}) {
	// catch panics caused by unexpected death of the server hosting the payload
	defer func() {
		recov = recover()
	}()

	// Access the container anchor of the currently-running payload of the virus.
	t := c.Walk(client.Split(payloadAnchor))
	// Wait until the payload process exits.
	t.Get().(client.Proc).Wait()
	// Remove the anchor of the now-dead payload process.
	// t.Scrub()

	// Wait a touch to prevent spinning, if the payload exits immediately every time it is run.
	time.Sleep(time.Second / 2)
	return
}

// The initial invocation of the virus:
//
//	virus DIALIN_CIRCUIT
//
// To invoke the virus in the role of a nucleus process:
//
//	virus DIALIN_CIRCUIT BACKCHAN_ANCHOR PAYLOAD_ANCHOR SELF_ANCHOR EPOCH
func main() {
	// Parse arguments
	var (
		err           error
		isNucleus     bool
		myAnchor      string
		payloadAnchor string
		epoch         int
	)
	switch len(os.Args) {
	case 2: // initial command-line invocation
	case 6: // invocation in role of nucleus
		isNucleus = true
		myAnchor = os.Args[4]
		payloadAnchor = os.Args[3]
		epoch, err = strconv.Atoi(os.Args[5])
		if err != nil {
			panic(err)
		}
	default:
		println("usage: virus circuit://...")
		os.Exit(1)
	}
	println("virus nucleus epoch", epoch, "dialing into", os.Args[1])
	c := client.Dial(os.Args[1], nil)

	// Create/get back channel
	backAnchor, backChan := findBackChan(c, isNucleus)

	// Acquire permission to send to back channel
	acquireBackChan(c, backChan, epoch)

	// The nucleus role waits for the payload process to die before it proceeds.
	if isNucleus {
		waitFotPayloadDeath(c, myAnchor, payloadAnchor, epoch)
	}

	payloadAnchor = spawnPayload(c, epoch)
	spawnNucleus(c, backAnchor, payloadAnchor, epoch)
}

// Create or get back channel
func findBackChan(c *client.Client, isNucleus bool) (backAnchor string, backChan client.Chan) {
	var err error
	if isNucleus {
		// The nucleus does not proceed with execution until it acquires permission
		// to send the the virus' back channel.
		backAnchor = os.Args[2]
		backChan = c.Walk(client.Split(backAnchor)).Get().(client.Chan)
	} else {
		// Make the back channel
		backServer := pickServer(c)
		//bc, err := backServer.Walk([]string{"virus", "back"}).Make(makers.ChanType, 3)
		backChan, err = client.MakeElement[client.Chan](backServer, 3, []string{"container-virus", "back"}...)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		backAnchor = path.Join("/", backServer.ServerID(), "virus", "back")
	}
	return
}

func acquireBackChan(c *client.Client, backChan client.Chan, epoch int) io.WriteCloser {
	// Acquire permission to send to back channel
	backPipe, err := backChan.Send()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(backPipe, "<%d> nucleus acquired back channel\n", epoch)
	return backPipe
}

func spawnPayload(c *client.Client, epoch int) (payloadAnchor string) {
	// Start the payload process
	// service := client.Cmd{
	// 	//Path: "/usr/bin/sayhe", // say is a standard OSX command which speaks, so it's easy to hear the virus in action.
	// 	Path:  "/usr/bin/echo", // say is a standard OSX command which speaks, so it's easy to hear the virus in action.
	// 	Args:  []string{"i am a virus"},
	// 	Scrub: true,
	// }
	opts := podman.ContainerCreateOptions{
		Image: "docker.io/library/redis:alpine",
	}

	// Randomly choose a circuit server to host the virus payload.
	a := pickServer(c)
	// Run the payload
	payloadEpoch := strconv.Itoa(epoch + 1)
	// p, err := a.Walk([]string{"virus", "payload", payloadEpoch}).Make(makers.ContainerType, opts)
	con, err := client.MakeElement[podman.Container](a, opts, []string{"container-virus", "payload", payloadEpoch}...)
	if err != nil {
		println("payload not created:", err.Error())
		os.Exit(1)
	}

	if _, err := con.Peek(); err != nil {
		println("payload not started:", err.Error())
		os.Exit(1)
	}

	// start the container
	if err = con.Start(); err != nil {
		println("container not started: ", err.Error())
		os.Exit(1)
	}

	// Close the standard input of the payload to indicate no intention to write data.
	con.Stdin().Close()

	//fmt.Fprintf(backPipe, "<%d> started payload\n", epoch)
	return path.Join("/", a.ServerID(), "container-virus", "payload", payloadEpoch)
}

func spawnNucleus(c *client.Client, backAnchor, payloadAnchor string, epoch int) {
	// Start the virus nucleus process, which will wait until the payload completes,
	// and then start a payload as well as a new nucleus elsewhere, over and over again.
	b := pickServer(c)
	virus, _ := filepath.Abs(os.Args[0]) // We assume that the virus binary is on the same path everywhere
	nucleusEpoch := strconv.Itoa(epoch + 1)
	nucleusAnchor := path.Join("/", b.ServerID(), "container-virus", "nucleus", nucleusEpoch)
	nucleus := client.Cmd{
		Path: virus,
		Args: []string{
			b.Addr(),      // dial-in circuit server address
			backAnchor,    // virus back channel anchor
			payloadAnchor, // payload anchor
			nucleusAnchor, // anchor of the spawned nucleus itself
			nucleusEpoch,  // epoch
		},
		Scrub: true,
	}

	p, err := client.MakeElement[client.Proc](b, nucleus, []string{"container-virus", "nucleus", nucleusEpoch}...)
	if err != nil {
		println("nucleus not created:", err.Error())
		os.Exit(1)
	}

	if err := p.Peek().Exit; err != nil {
		println("nucleus not started:", err.Error())
		os.Exit(1)
	}

	p.Stdin().Close()
}
