package wasm

import "io"

type Wasm interface {
	Scrub()
	IsDone() bool
	Peek() (*Status, error)
	PeekBytes() []byte
	Signal(sig string) error
	Wait() (*Status, error)
	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}
