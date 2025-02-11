package makers

import (
	"fmt"
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/proc"
)

func init() {
	client.RegisterElementMaker(&procElementMaker{
		client.NewBaseElementMaker("proc", reflect.TypeOf(proc.YProc{})),
	})
}

var ProcType = reflect.TypeOf(proc.YProc{})

// implementation for a specific maker
type procElementMaker struct {
	client.BaseElementMaker
}

func (p *procElementMaker) Make(y anchor.YTerminal, arg any) (v any, err error) {
	v, err = y.Make(p.Name(), arg)
	if err != nil {
		return nil, err
	}

	// Check type of v
	if reflect.TypeOf(v) != p.Type() {
		return nil, fmt.Errorf("client/circuit mismatch, kind=%v", p.Name())
	}

	// return value can now be type asserted to t whithout error
	return yprocProc{v.(proc.YProc)}, nil
}

func (p *procElementMaker) Get(v any) any {
	return yprocProc{v.(proc.YProc)}
}

func statstat(s proc.Stat) client.ProcStat {
	return client.ProcStat{
		Cmd:   retypeProcStat(s.Cmd),
		Exit:  s.Exit,
		Phase: s.Phase,
	}
}

func retypeProcStat(c proc.Cmd) client.Cmd {
	return client.Cmd{
		Env:   c.Env,
		Dir:   c.Dir,
		Path:  c.Path,
		Args:  c.Args,
		Scrub: c.Scrub,
	}
}

type yprocProc struct {
	proc.YProc
}

func (y yprocProc) Wait() (client.ProcStat, error) {
	s, err := y.YProc.Wait()
	if err != nil {
		return client.ProcStat{}, err
	}
	return statstat(s), nil
}

func (y yprocProc) GetCmd() client.Cmd {
	return retypeProcStat(y.YProc.GetCmd())
}

func (y yprocProc) Peek() client.ProcStat {
	return statstat(y.YProc.Peek())
}
