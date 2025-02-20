// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2014 Petar Maymounkov <p@gocircuit.org>

package dns

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/miekg/dns"
)

type Nameserver interface {
	client.Nameserver
	X() circuit.X
}

type nameserver struct {
	sync.Mutex
	server *dns.Server
	addr   net.Addr
	rr     map[string][]dns.RR // name -> rr
}

func init() {
	anchor.RegisterElement("dns", ef, yf)
}

func MakeNameserver(addr string) (_ Nameserver, err error) {
	ns := &nameserver{
		rr: make(map[string][]dns.RR),
	}
	if err = ns.startUdpServer(addr); err != nil {
		return nil, err
	}
	return ns, nil
}

func (ns *nameserver) startUdpServer(addr string) error {
	pc, err := net.ListenPacket("udp", addr) // empty-string address picks an available port on 0.0.0.0
	if err != nil {
		return err
	}
	ns.server = &dns.Server{
		PacketConn: pc,
		Handler:    ns,
	}
	ns.addr = pc.LocalAddr()
	go func() {
		ns.server.ActivateAndServe()
		pc.Close()
	}()
	return nil
}

func (ns *nameserver) lookup(q string) []dns.RR {
	ns.Lock()
	defer ns.Unlock()
	return ns.rr[q]
}

func (ns *nameserver) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(req)
	q := msg.Question[0].Name // question, e.g. "miek.nl."

	rr := ns.lookup(q)
	if len(rr) == 0 { // no entry
		w.Close()
		return
	}

	msg.Answer = make([]dns.RR, 1)
	msg.Answer[0] = rr[0]
	w.WriteMsg(msg)
}

func (ns *nameserver) Scrub() {
	ns.Lock()
	defer ns.Unlock()
	if ns.server == nil {
		return
	}
	ns.server.Shutdown()
	ns.server = nil
}

func (ns *nameserver) X() circuit.X {
	return circuit.Ref(XNameserver{ns})
}

func (ns *nameserver) Set(rr string) error {
	ss, err := dns.NewRR(rr)
	if err != nil {
		return err
	}
	ns.Lock()
	defer ns.Unlock()
	ns.rr[ss.Header().Name] = append(ns.rr[ss.Header().Name], ss)
	return nil
}

func (ns *nameserver) Unset(name string) {
	ns.Lock()
	defer ns.Unlock()
	delete(ns.rr, name)
}

func (ns *nameserver) Peek() client.NameserverStat {
	ns.Lock()
	defer ns.Unlock()
	var stat client.NameserverStat
	stat.Address = ns.addr.String()
	stat.Records = make(map[string][]string)
	for name, rr := range ns.rr {
		var ss []string
		for _, record := range rr {
			ss = append(ss, record.String())
		}
		stat.Records[name] = ss
	}
	return stat
}

func (ns *nameserver) PeekBytes() []byte {
	b, _ := json.MarshalIndent(ns.Peek(), "", "\t")
	return b
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	ns, err := MakeNameserver(arg.(string))
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func yf(x circuit.X) (any, error) {
	return YNameserver{x}, nil
}
