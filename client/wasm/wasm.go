package wasm

import "io"

type Wasm interface {
	Scrub()
	IsDone() bool
	Peek() (*Status, error)
	Signal(sig string) error
	Wait() (*Status, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}
