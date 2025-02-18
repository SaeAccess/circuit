// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		// stdin, stdout, stderr
		{
			Name:      "stdin",
			Usage:     "Forward this tool's standard input to that of the process",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    stdin,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "stdout",
			Usage:     "Forward the standard output of the process to the standard output of this tool",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    stdout,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
		{
			Name:      "stderr",
			Usage:     "Forward the standard error of the process to the standard output of this tool",
			Args:      true,
			ArgsUsage: "Anchor",
			Action:    stderr,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

func stdin(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("stdin needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(interface {
		Stdin() io.WriteCloser
	})
	if !ok {
		return errors.New("anchor does not implement the 'Stdin() io.WriteCloser' interface method. See Proc or container")
	}
	q := u.Stdin()
	if _, err = io.Copy(q, os.Stdin); err != nil {
		return errors.Wrapf(err, "transmission error: %v", err)
	}
	if err = q.Close(); err != nil {
		return errors.Wrapf(err, "error closing stdin: %v", err)
	}
	return
}

func stdout(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("stdout needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(interface {
		Stdout() io.ReadCloser
	})
	if !ok {
		return errors.New("anchor does not implement the 'Stdout() io.ReadCloser' interface method. See Proc or container")
	}
	io.Copy(os.Stdout, u.Stdout())
	return
}

func stderr(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return errors.New("stderr needs one anchor argument")
	}
	w, _ := parseGlob(args.First())
	u, ok := c.Walk(w).Get().(interface {
		Stderr() io.ReadCloser
	})
	if !ok {
		return errors.New("anchor does not implement the 'Stderr() io.ReadCloser' interface method. See Proc or container")
	}
	io.Copy(os.Stdout, u.Stderr())
	// if _, err := io.Copy(os.Stdout, u.Stderr()); err != nil {
	// 	fatalf("transmission error: %v", err)
	// }
	return
}
