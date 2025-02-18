// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package valve

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/use/circuit"
)

type Valve interface {
	client.Chan
	IsDone() bool
	X() circuit.X
}

// valve
type valve struct {
	send struct {
		abr <-chan struct{} // abort when closed
		sync.Mutex
		tun chan<- interface{}
	}
	recv struct {
		abr <-chan struct{} // abort when closed
		tun <-chan interface{}
	}
	ctrl struct {
		sync.Mutex
		abr  chan<- struct{}
		stat client.ChanStat
	}
}

// type Stat struct {
// 	client.ChanStat
// 	Opened bool `json:"opened,omitempty"`
// }

func init() {
	anchor.RegisterElement("chan", ef, yf)
}

// Sender-receiver pipe capacity (once matched)
const MessageCap = 32e3 // 32K

func MakeValve(n int) Valve {
	v := &valve{}
	tun, abr := make(chan interface{}, n), make(chan struct{})
	v.send.tun, v.recv.tun = tun, tun
	v.ctrl.abr, v.send.abr, v.recv.abr = abr, abr, abr
	v.ctrl.stat.Opened, v.ctrl.stat.Cap = true, n
	return v
}

func (v *valve) X() circuit.X {
	return circuit.Ref(XValve{v})
}

func (v *valve) incSend() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumSend++
}

func (v *valve) incRecv() {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	v.ctrl.stat.NumRecv++
}

// Cap returns the capacity of the valve and whether it was set.
func (v *valve) Cap() int {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	if v.ctrl.stat.Opened {
		return v.ctrl.stat.Cap
	}
	return -1
}

func (v *valve) Stat() client.ChanStat {
	v.ctrl.Lock()
	defer v.ctrl.Unlock()
	return v.ctrl.stat
}

func (v *valve) PeekBytes() []byte {
	b, _ := json.MarshalIndent(v.Stat(), "", "\t")
	return b
}

type scrubValve struct {
	t *anchor.Terminal
	Valve
}

func (v *scrubValve) Close() error {
	defer func() {
		if v.Valve.IsDone() {
			v.t.Scrub()
		}
	}()
	return v.Valve.Close()
}

func (v *scrubValve) Recv() (io.ReadCloser, error) {
	defer func() {
		if v.Valve.IsDone() {
			v.t.Scrub()
		}
	}()
	return v.Valve.Recv()
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	capacity, ok := arg.(int)
	if !ok {
		return nil, fmt.Errorf("invalid argument to chan element factory, arg=%T", arg)
	}

	return &scrubValve{t, MakeValve(capacity)}, nil
}

func yf(x circuit.X) (any, error) {
	return YValve{X: x}, nil
}
