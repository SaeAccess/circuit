// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/gocircuit/circuit/element/podman/container"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		{
			Name:   "waitall",
			Usage:  "Wait until a set of processes all exit",
			Action: waitall,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// TODO use Waiter interface to determine the anchors that can wait
func waitall(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error: %v", r)
		}
	}()

	c := dial(x)
	s := bufio.NewScanner(os.Stdin)
	var t []string
	ch := make(chan int)
	for s.Scan() {
		src := s.Text()
		i := len(t)
		t = append(t, src)
		go func() { // wait on w
			w, _ := parseGlob(src)
			var e error
			switch u := c.Walk(w).Get().(type) {
			case client.Proc:
				_, e = u.Wait()
			case docker.Container:
				_, e = u.Wait()
			case container.Container:
				_, e = u.Wait()
			default:
				println("anchor", w, " is not a process or a container")
			}
			if e != nil {
				log.Fatal(errors.Errorf("wait error: %v", e))
			}
			ch <- i
		}()
	}
	for range t {
		<-ch
	}
	return
}
