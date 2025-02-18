// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gocircuit/circuit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func fatalf(format string, arg ...interface{}) {
	println(fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func readkey(x *cli.Context) (key []byte) {
	var hmac string
	if hmac = x.String("hmac"); !x.IsSet("hmac") {
		return nil
	}

	b64, err := os.ReadFile(hmac)
	if err != nil {
		fatalf("problem reading private key file (%s): %v", hmac, err)
	}
	if key, err = base64.StdEncoding.DecodeString(string(b64)); err != nil {
		fatalf("problem decoding base64 private key: %v", err)
	}
	return
}

func dial(x *cli.Context) *client.Client {
	switch {
	case x.String("dial") != "":
		defer func() {
			if r := recover(); r != nil {
				fatalf("addressed server is gone or authentication failed")
			}
		}()
		return client.Dial(x.String("dial"), readkey(x))

	case os.Getenv("CIRCUIT") != "":
		//buf, err := ioutil.ReadFile(os.Getenv("CIRCUIT"))
		buf, err := os.ReadFile(os.Getenv("CIRCUIT"))
		if err != nil {
			fatalf("circuit environment file %s is not readable: %v", os.Getenv("CIRCUIT"), err)
		}
		defer func() {
			if r := recover(); r != nil {
				fatalf("addressed server is gone or authentication failed")
			}
		}()
		log.Printf("dial string from env: %s", string(buf))
		return client.Dial(strings.TrimSpace(string(buf)), readkey(x))

	case x.String("discover") != "":
		// try multicast as a last
		defer func() {
			if r := recover(); r != nil {
				fatalf("multicast address is unresponsive or authentication failed")
			}
		}()
		return client.DialDiscover(x.String("discover"), readkey(x))

	}
	fatalf("no dial or discovery addresses available; use -dial or -discover")
	panic(0)
}

func readStdin[T any]() (*T, error) {
	buf, _ := io.ReadAll(os.Stdin)
	var opts T
	if err := json.Unmarshal(buf, &opts); err != nil {
		return nil, errors.Wrapf(err, "failed to parse json options for command from stdin: %v", err)
	}

	return &opts, nil
}

func ParseAnchor(x *cli.Context) (a []string, err error) {
	args := x.Args()
	if args.Len() < 1 {
		return nil, errors.New("anchor argument is required")
	}

	w, _ := parseGlob(args.First())
	return w, nil
}

type EmptyInterface interface{}

func getAnchorType[T any](x *cli.Context, t string) (_ T, err error) {
	// connect
	c := dial(x)

	var w []string
	if w, err = ParseAnchor(x); err != nil {
		var z T
		return z, err
	}

	vol, ok := c.Walk(w).Get().(T)
	if !ok {
		var z T
		return z, fmt.Errorf("anchor is not a podman %s", t)
	}

	return vol, nil
}

type Exists interface {
	Exists() error
}

// TODO should return bool
func checkAnchorType[T Exists](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var at T
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	if err = at.Exists(); err != nil {
		return errors.Wrapf(err, "%s does not exist", t)
	}

	return
}

type Inspecter[D any] interface {
	Inspect() (*D, error)
}

func inspectAnchorType[T Inspecter[D], D any](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var at T
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	var data *D
	if data, err = at.Inspect(); err != nil {
		return errors.Wrapf(err, "podman %s inspect error", t)
	}

	// write to stdout
	buf, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(buf))
	return
}

type Reloader interface {
	Reload() error
}

func reloadAnchorType[T Reloader](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var at Reloader
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	if err = at.Reload(); err != nil {
		return errors.Wrapf(err, "podman %s reload error", t)
	}

	return
}

type Updater[OPT any] interface {
	Update(opt *OPT) error
}

func updateAnchorType[T Updater[OPT], OPT any](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var at Updater[OPT]
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	var opts *OPT
	if opts, err = readStdin[OPT](); err != nil {
		return err
	}

	if err = at.Update(opts); err != nil {
		return errors.Wrapf(err, "podman %s reload error", t)
	}

	return
}

type Pauser interface {
	Pause() error
}

func pauseAnchorType[T Pauser](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var at Pauser
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	if err = at.Pause(); err != nil {
		return errors.Wrapf(err, "podman %s pause error", t)
	}

	return
}

type Restarter interface {
	Restart() error
}

func restartAnchorType[T Restarter](x *cli.Context, t string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var at T
	if at, err = getAnchorType[T](x, t); err != nil {
		return
	}

	if err = at.Restart(); err != nil {
		return errors.Wrapf(err, "podman %s restart error", t)
	}

	return
}
