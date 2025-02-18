package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/proc"
	"github.com/pkg/errors"
)

var ProcType = reflect.TypeOf((*client.Proc)(nil)).Elem()

func init() {
	client.RegisterElementMaker(&procElementMaker{
		client.NewBaseElementMaker("proc", ProcType),
	})
}

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
	if !reflect.TypeOf(v).Implements(p.Type()) {
		return nil, errors.Wrapf(client.ErrMismatchType, "%v does not implement %v", reflect.TypeOf(v), ProcType)
	}

	// return value can now be type asserted to t whithout error
	return yprocProc{v.(proc.YProc)}, nil
}

// func statstat(s proc.Stat) client.ProcStat {
// 	return client.ProcStat{
// 		Cmd:   retypeProcStat(s.Cmd),
// 		Exit:  s.Exit,
// 		Phase: s.Phase,
// 	}
// }

// func retypeProcStat(c proc.Cmd) client.Cmd {
// 	return client.Cmd{
// 		Env:   c.Env,
// 		Dir:   c.Dir,
// 		Path:  c.Path,
// 		Args:  c.Args,
// 		Scrub: c.Scrub,
// 	}
// }

type yprocProc struct {
	proc.YProc
}

func (y yprocProc) Wait() (client.ProcStat, error) {
	s, err := y.YProc.Wait()
	if err != nil {
		return client.ProcStat{}, err
	}
	return s, nil
}

func (y yprocProc) GetCmd() client.Cmd {
	return y.YProc.GetCmd()
}

func (y yprocProc) Peek() client.ProcStat {
	return y.YProc.Peek()
}
