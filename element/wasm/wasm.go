package wasm

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/gocircuit/circuit/anchor"
	ws "github.com/gocircuit/circuit/client/wasm"
	"github.com/gocircuit/circuit/element/proc"
	"github.com/gocircuit/circuit/use/circuit"
	"github.com/second-state/WasmEdge-go/wasmedge"
)

// Statistics wasm runtime statistics
type Statistics struct {
}

type Wasm interface {
	ws.Wasm
	X() circuit.X
}

type container struct {
	name string
	cmd  *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	exit   <-chan error

	vm *wasmedge.VM
}

func init() {
	anchor.RegisterElement("wasm", ef, yf)
}

func MakeContainer(exe *ws.Execute) (_ Wasm, err error) {
	wasmedge.SetLogErrorLevel()

	// build wasm config from execute config
	conf := buildConfigure(exe)

	//var conf = wasmedge.NewConfigure(wasmedge.REFERENCE_TYPES)
	conf.AddConfig(wasmedge.WASI)
	var vm = wasmedge.NewVMWithConfig(conf)
	var wasi = vm.GetImportModule(wasmedge.WASI)

	ch := make(chan error, 1)
	con := &container{}

	wasi.InitWasi(
		os.Args[1:],     // The args
		os.Environ(),    // The envs
		[]string{".:."}, // The mapping directories
	)

	go func() {
		ch <- con.cmd.Wait()
		close(ch)
		con.cmd.Stdout.(io.Closer).Close()
		con.cmd.Stderr.(io.Closer).Close()
	}()

	// Determine if functinal or application wasm execution
	// Instantiate wasm. _start refers to the main() function
	vm.RunWasmFile(os.Args[1], "_start")

	runtime.SetFinalizer(con,
		func(c *container) {
			c.vm.Release()
		},
	)

	conf.Release()

	return con, nil
}

func buildConfigure(exe *ws.Execute) *wasmedge.Configure {
	conf := &wasmedge.Configure{}

	return conf
}

func (con *container) Wait() (_ *ws.Status, err error) {
	<-con.exit
	return con.Peek()
}

func (con *container) Stdin() io.WriteCloser {
	return con.stdin
}

func (con *container) Stdout() io.ReadCloser {
	return con.stdout
}

func (con *container) Stderr() io.ReadCloser {
	return con.stderr
}

func (con *container) Peek() (stat *ws.Status, err error) {
	return
}

func (con *container) PeekBytes() []byte {
	s, err := con.Peek()
	if err != nil {
		return []byte{}
	}

	b, _ := json.MarshalIndent(s, "", "\t")
	return b
}

func (con *container) Scrub() {
	//exec.Command(dkr, "rm", con.name).Run()
}

func (con *container) Signal(sig string) error {
	signo, ok := proc.ParseSignal(sig)
	if !ok {
		return errors.New("signal name not recognized")
	}
	if con.cmd.Process == nil {
		return errors.New("no process")
	}
	return con.cmd.Process.Signal(signo)
}

func (con *container) IsDone() bool {
	select {
	case <-con.exit:
		return true
	default:
		return false
	}
}

func (con *container) X() circuit.X {
	return circuit.Ref(XWasm{con})
}

func ef(t *anchor.Terminal, arg any) (anchor.Element, error) {
	exe, ok := arg.(ws.Execute)
	if !ok {
		return nil, errors.New("invalid wasm argument")
	}

	w, err := MakeContainer(&exe)
	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			recover()
		}()

		if exe.Scrub {
			defer t.Scrub()
		}

		w.Wait()
	}()

	return w, nil
}

func yf(x circuit.X) (any, error) {
	return YWasm{x}, nil
}
