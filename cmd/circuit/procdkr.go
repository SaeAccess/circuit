// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/docker"
	"github.com/gocircuit/circuit/client/makers"
	"github.com/gocircuit/circuit/element/podman/container"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		// proc/dkr-specific
		{
			Name:      "mkdkr",
			Usage:     "Create a docker container element",
			Args:      true,
			ArgsUsage: "anchor",
			// TODO add Before to check for docker
			Action: mkdkr,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "mkproc",
			Usage:     "Create a process element",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    mkproc,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "signal",
			Usage:     "Send a signal to a running process",
			Args:      true,
			ArgsUsage: "anchor signal",
			Action:    sgnl,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "wait",
			Usage:     "Wait until a process exits",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    wait,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// circuit mkproc /X1234/hola/charlie << EOF
// { … }
// EOF
// TODO: Proc element disappears if command misspelled and error condition not obvious.
func mkproc(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("mkproc needs an anchor argument")
	}
	w, _ := parseGlob(args.First())
	buf, _ := io.ReadAll(os.Stdin)
	var cmd client.Cmd
	if err = json.Unmarshal(buf, &cmd); err != nil {
		return errors.Wrapf(err, "command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		cmd.Scrub = true
	}
	p, err := c.Walk(w).Make(makers.ProcType, cmd)
	if err != nil {
		return errors.Wrapf(err, "mkproc error: %s", err)
	}

	if proc, ok := p.(client.Proc); ok {
		ps := proc.Peek()
		if ps.Exit != nil {
			return errors.Errorf("%v", ps.Exit)
		}
	}

	return
}

func mkdkr(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("mkdkr needs an anchor argument")
	}
	w, _ := parseGlob(args.First())
	buf, _ := io.ReadAll(os.Stdin)
	var run docker.Run
	if err = json.Unmarshal(buf, &run); err != nil {
		return errors.Wrapf(err, "command json not parsing: %v", err)
	}
	if x.Bool("scrub") {
		run.Scrub = true
	}

	// get client side anchor to create the docker container
	if _, err = c.Walk(w).Make(makers.DockerType, run); err != nil {
		return errors.Wrapf(err, "mkdkr error: %s", err)
	}
	return
}

// circuit signal kill /X1234/hola/charlie
func sgnl(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 2 {
		return errors.New("signal needs an anchor and a signal name arguments")
	}
	w, _ := parseGlob(args.Get(1))
	u, ok := c.Walk(w).Get().(interface {
		Signal(string) error
	})
	if !ok {
		return errors.New("anchor is not a process or a docker container")
	}
	if err = u.Signal(args.First()); err != nil {
		return errors.Wrapf(err, "signal error: %v", err)
	}
	return
}

// TODO implement Waiter interface
func wait(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("wait needs one anchor argument")
	}
	w, _ := parseGlob(args.First())

	var stat interface{}
	switch u := c.Walk(w).Get().(type) {
	case client.Proc:
		stat, err = u.Wait()
	case docker.Container:
		stat, err = u.Wait()
	case container.Container:
		stat, err = u.Wait()
	default:
		return errors.New("anchor is not a process or a container")
	}
	if err != nil {
		return errors.Wrapf(err, "wait error: %v", err)
	}
	buf, _ := json.MarshalIndent(stat, "", "\t")
	fmt.Println(string(buf))
	return
}
