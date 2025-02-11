// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/client/makers"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

func init() {
	// nameserver
	var cmds = []*cli.Command{
		{
			Name:      "mkdns",
			Usage:     "Create a nameserver element",
			Args:      true,
			ArgsUsage: "anchor [address]",
			Action:    mkdns,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "set",
			Usage:     "Set a resource record in a nameserver element",
			Args:      true,
			ArgsUsage: "anchor resource-record",
			Action:    nset,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "unset",
			Usage:     "Remove all resource records for a name in a nameserver element",
			Args:      true,
			ArgsUsage: "anchor resource-name",
			Action:    nunset,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

func mkdns(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() < 1 {
		return errors.New("mkdns needs an anchor and an optional address arguments")
	}
	var addr string
	if args.Len() == 2 {
		addr = args.Get(1)
	}
	w, _ := parseGlob(args.First())

	//if _, err = c.Walk(w).MakeNameserver(addr); err != nil {
	if _, err = c.Walk(w).Make(makers.NameserverType, addr); err != nil {
		return errors.Wrapf(err, "mkdns error: %s", err)
	}
	return
}

func nset(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 2 {
		return errors.New("set needs an anchor and a resource record arguments")
	}
	w, _ := parseGlob(args.First())
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		err := u.Set(args.Get(1))
		if err != nil {
			return errors.Wrapf(err, "set resoure record error: %v", err)
		}
	default:
		return errors.New("not a nameserver element")
	}
	return
}

func nunset(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 2 {
		return errors.New("unset needs an anchor and a resource name arguments")
	}
	w, _ := parseGlob(args.First())
	switch u := c.Walk(w).Get().(type) {
	case client.Nameserver:
		u.Unset(args.Get(1))
	default:
		return errors.New("not a nameserver element")
	}
	return
}
