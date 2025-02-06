// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/use/n"
	"github.com/pkg/errors"

	// "bytes"
	"io"
	"os"

	"github.com/gocircuit/circuit/client"

	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		// server-specific
		{
			Name:      "stk",
			Usage:     "Print the runtime stack trace of a server element",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    stack,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "join",
			Usage:     "Merge the networks of this circuit server and that of the argument circuit address",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    join,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "suicide",
			Usage:     "Kill a chosen circuit daemon",
			Args:      true,
			ArgsUsage: "anchor circuit-address",
			Action:    stack,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}
	RegisterCommand(cmds...)
}

func stack(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("recv needs one anchor argument")
	}

	w, _ := parseGlob(args.First())
	switch u := c.Walk(w).Get().(type) {
	case client.Server:
		r, err := u.Profile("goroutine")
		if err != nil {
			return errors.Wrapf(err, "error: %v", err)
		}
		io.Copy(os.Stdout, r)
	default:
		return errors.New("not a server")
	}
	return
}

func suicide(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("suicide needs one server anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(client.Server)
	if !ok {
		return errors.New("not a server")
	}
	u.Suicide()
	return
}

func join(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()
	c := dial(x)
	args := x.Args()
	if args.Len() != 2 {
		return errors.New("join needs one anchor argument and one circuit address argument")
	}
	// Verify the target circuit address is valid
	if _, err = n.ParseAddr(args.Get(1)); err != nil {
		return errors.Wrapf(err, "argument %q is not a valid circuit address", args.Get(1))
	}
	//
	w, _ := parseGlob(args.First())
	switch u := c.Walk(w).Get().(type) {
	case client.Server:
		if err = u.Rejoin(args.Get(1)); err != nil {
			return errors.Wrapf(err, "error: %v", err)
		}
	default:
		return errors.New("not a server")
	}
	return
}
