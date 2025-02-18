// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"fmt"

	"github.com/gocircuit/circuit/anchor"
	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		{
			Name:      "peek",
			Usage:     "Query element state asynchronously",
			Args:      true,
			ArgsUsage: "anchor",
			Action:    peek,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:   "scrub",
			Usage:  "Abort and remove an element",
			Action: scrb,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// circuit peek /X1234/hola/charlie
// TODO change this to resolve Peeker/Inspector interface on the Anchor so we can type switch on it
func peek(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("peek needs one anchor argument")
	}

	w, _ := parseGlob(args.First())

	// Check if anchor supports the Peeker interface
	av := c.Walk(w).Get()
	if p, ok := av.(anchor.Peeker); ok {
		fmt.Println(string(p.PeekBytes()))
	} else {
		return fmt.Errorf("anchor type '%T' does not support the anchor.Peeker interface", av)
	}

	return
}

func scrb(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("scrub needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	c.Walk(w).Scrub()
	return
}
