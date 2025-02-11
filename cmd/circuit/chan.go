// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	//"log"
	"io"
	"os"
	"strconv"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/makers"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

func init() {
	// channel-specific commands, should make into subcommands
	cmds := []*cli.Command{
		{
			Name:      "mkchan",
			Usage:     "Create a channel element",
			Args:      true,
			ArgsUsage: "anchor capacity",
			Action:    mkchan,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "send",
			Usage:     "Send data to the channel from standard input",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    send,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "recv",
			Usage:     "Receive data from a channel or a subscription on stadard output",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    recv,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "close",
			Usage:     "Close the channel after all current transmissions complete",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    clos,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// circuit mkchan /X1234/hola/charlie 0
func mkchan(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 2 {
		return errors.New("mkchan needs an anchor and a capacity arguments")
	}
	w, _ := parseGlob(args.First())
	a := c.Walk(w)
	n, err := strconv.Atoi(args.Get(1))
	if err != nil || n < 0 {
		return errors.New("second argument to mkchan must be a non-negative integral capacity")
	}
	if _, err := a.Make(makers.ChanType, n); err != nil {
		// if _, err = a.MakeChan(n); err != nil {
		return errors.Wrapf(err, "mkchan error: %s", err)
	}
	return
}

// TODO add a file argument to send, and if specified dont read from os.Stdin
func send(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("send needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		return errors.New("not a channel")
	}
	msgw, err := u.Send()
	if err != nil {
		return errors.Wrapf(err, "send error: %v", err)
	}
	if _, err = io.Copy(msgw, os.Stdin); err != nil {
		return errors.Wrapf(err, "transmission error: %v", err)
	}
	return
}

func clos(x *cli.Context) (err error) {
	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("close needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(client.Chan)
	if !ok {
		return errors.New("not a channel")
	}
	if err := u.Close(); err != nil {
		return errors.Wrapf(err, "close error: %v", err)
	}
	return
}
