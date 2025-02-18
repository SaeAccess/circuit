// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package cmd

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/urfave/cli/v2"
)

func init() {
	var cmds = []*cli.Command{
		{
			Name:   "keygen",
			Usage:  "Generate a new random HMAC key",
			Action: keygen,
		},
	}

	RegisterCommand(cmds...)
}

func keygen(c *cli.Context) (err error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//rand.Seed(time.Now().UnixNano())
	seed := make([]byte, 32)
	for i := range seed {
		//seed[i] = byte(rand.Int31())
		seed[i] = byte(r.Int31())
	}
	key := sha512.Sum512(seed)
	text := base64.StdEncoding.EncodeToString(key[:])
	fmt.Println(text)
	return
}
