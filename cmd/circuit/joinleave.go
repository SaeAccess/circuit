// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

import (
	"github.com/gocircuit/circuit/client/makers"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		// subscription-specific
		{
			Name:      "mk@join",
			Usage:     "Create a subscription element, receiving server join events",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    mkonjoin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "mk@leave",
			Usage:     "Create a subscription element, receiving server leave events",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    mkonleave,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// circuit mk@join /X1234/hola/listy
func mkonjoin(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("mk@join needs an anchor argument")
	}
	w, _ := parseGlob(args.First())
	//if _, err = c.Walk(w).MakeOnJoin(); err != nil {
	if _, err = c.Walk(w).Make(makers.JoinType, ""); err != nil {
		return errors.Wrapf(err, "mk@join error: %s", err)
	}
	return
}

// circuit mk@leave /X1234/hola/listy
func mkonleave(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("mk@leave needs an anchor argument")
	}
	w, _ := parseGlob(args.First())
	//if _, err = c.Walk(w).MakeOnLeave(); err != nil {
	if _, err = c.Walk(w).Make(makers.LeaveType, ""); err != nil {
		return errors.Wrapf(err, "mk@leave error: %s", err)
	}
	return
}
