// Copyright 2013 The Go Circuit Project
// Use of this source code is governed by the license for
// The Go Circuit Project, found in the LICENSE file.
//
// Authors:
//   2013 Petar Maymounkov <p@gocircuit.org>

package client

import (
	"encoding/json"
	"io"
)

// Proc provides access to a circuit process element.
// All methods panic if the hosting circuit server dies.
type Proc interface {

	// Wait blocks until the underlying OS process exits and returns the final status of the process.
	// An error is returned only if the wait invocation is aborted by a concurring call to Scrub.
	Wait() (ProcStat, error)

	// Signal sends an OS signal to the process. The following are recognized signal names:
	// ABRT, ALRM, BUS, CHLD, CONT, FPE, HUP, ILL, INT, IO, IOT,  KILL, PIPE,
	// PROF, QUIT, SEGV,  STOP, SYS, TERM, TRAP, TSTP, TTIN, TTOU,  URG, USR1,
	// USR2, VTALRM, WINCH, XCPU, XFSZ.
	Signal(sig string) error

	// GetEnv returns the environment at the hosting server OS.
	GetEnv() []string

	// GetCmd returns the command that started this process.
	GetCmd() Cmd

	// Peek asynchronously returns the current state of the process.
	Peek() ProcStat

	PeekBytes() []byte

	// Scrub abandons the circuit process element, without affecting the underlying OS process.
	Scrub()

	// Stdin returns a WriterCloser to the standard input of the underlying OS process.
	// The user is responsible for closing the standard input, even if they do not
	// intend to write to it.
	Stdin() io.WriteCloser

	// Stdout returns the standard output of the underlying OS process.
	Stdout() io.ReadCloser

	// Stderr returns the standard error of the underlying OS process.
	Stderr() io.ReadCloser
}

// Cmd describes the execution parameters for an OS process.
type Cmd struct {

	// Env, if set, is the desired OS execution environment. It corresponds to Cmd.Env from package "os/exec".
	Env []string `json:"env,omitempty"`

	// Dir, if non-empty, is the working directory for the process.
	Dir string `json:"dir,omitempty"`

	// Path is the local file-system path, at the respective circuit server, to the process binary.
	Path string `json:"path,omitempty"`

	// Args is a list of command line arguments to be passed on to the process.
	// The first element in the slice corresponds to the first argument to the process (not to its binary path).
	Args []string `json:"args,omitempty"`

	// If Scrub is set, the process element will automatically be removed from its anchor
	// when the process exits.
	Scrub bool `json:"scrub,omitempty"`
}

func ParseCmd(src string) (*Cmd, error) {
	x := &Cmd{}
	if err := json.Unmarshal([]byte(src), x); err != nil {
		return nil, err
	}
	return x, nil
}

func (x Cmd) String() string {
	b, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(0)
	}
	return string(b)
}

// func (cmd Cmd) retype() proc.Cmd {
// 	return proc.Cmd{
// 		Env:   cmd.Env,
// 		Dir:   cmd.Dir,
// 		Path:  cmd.Path,
// 		Args:  cmd.Args,
// 		Scrub: cmd.Scrub,
// 	}
// }

// ProcStat encloses process state information.
type ProcStat struct {

	// Cmd is a copy of the command that started the process.
	Cmd Cmd `json:"cmd,omitempty"`

	// Error will be non-nil if the process has already exited in error.
	Exit error `json:"exit,omitempty"`

	// Phase describes the current state of the process.
	// Its possible values are Running, Exited, Stopped, Signaled, Continued and Unknown.
	Phase string `json:"phase,omitempty"`
}

const (
	Running   = "running"
	Exited    = "exited"
	Stopped   = "stopped"
	Signaled  = "signaled"
	Continued = "continued"
	Unknown   = "unknown"
)
