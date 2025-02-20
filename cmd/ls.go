// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		{
			Name:      "ls",
			Usage:     "List circuit elements",
			Args:      true,
			ArgsUsage: "glob",
			Action:    ls,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
				&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
				&cli.BoolFlag{Name: "long", Aliases: []string{"l"}, Usage: "show detailed anchor information"},
				&cli.BoolFlag{Name: "depth", Aliases: []string{"de"}, Usage: "traverse anchors in depth-first order (leaves first)"},
				&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
			},
		},
	}

	RegisterCommand(cmds...)
}

// circuit ls /Q123/apps/charlie
// circuit ls /...
func ls(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			if err, ok = r.(error); !ok {
				err = errors.Errorf("%v", r)
			} else {
				err = errors.Wrapf(r.(error), "error, likely due to missing server or misspelled anchor: %v", r)
			}
		}
	}()
	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		println("ls needs a glob argument")
		os.Exit(1)
	}
	w, ellipses := parseGlob(args.First())
	list(0, "/", c.Walk(w), ellipses, x.Bool("long"), x.Bool("depth"))
	return
}

func list(level int, prefix string, anchor client.Anchor, recurse, long, depth bool) {
	if anchor == nil {
		return
	}

	//println(fmt.Sprintf("prefix=%v a=%v/%T r=%v", prefix, anchor, anchor, recurse))
	var c children
	for n, a := range anchor.View() {
		e := &entry{n: n, a: a}
		if v := a.Get(); v == nil {
			e.k = "."
		} else {
			if maker := client.FindElementMaker(v); maker != nil {
				e.k = maker.Name()
			} else {
				e.k = "."
			}
		}
		c = append(c, e)
	}
	sort.Sort(c)
	for _, e := range c {
		if recurse && depth {
			list(level+1, prefix+e.n+"/", e.a, true, long, depth)
		}
		if long {
			fmt.Printf("%-15s %s%s\n", e.k, prefix, e.n)
		} else {
			fmt.Printf("%s%s\n", prefix, e.n)
		}
		if recurse && !depth {
			list(level+1, prefix+e.n+"/", e.a, true, long, depth)
		}
	}
}

type entry struct {
	n string
	a client.Anchor
	k string
}

type children []*entry

func (c children) Len() int {
	return len(c)
}

func (c children) Less(i, j int) bool {
	return c[i].n < c[j].n
}

func (c children) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func parseGlob(pattern string) (walk []string, ellipses bool) {
	for _, p := range strings.Split(pattern, "/") {
		if len(p) == 0 {
			continue
		}
		walk = append(walk, p)
	}
	if len(walk) == 0 {
		return
	}
	if walk[len(walk)-1] == "..." {
		walk = walk[:len(walk)-1]
		ellipses = true
	}
	return
}
