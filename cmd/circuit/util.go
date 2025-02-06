// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package main

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
