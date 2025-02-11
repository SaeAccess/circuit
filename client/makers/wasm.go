package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/wasm"
)

func init() {
	client.RegisterElementMaker(&wasmElementMaker{
		client.NewBaseElementMaker("wasm", reflect.TypeOf(wasm.YWasm{})),
	})
}

// implementation for a specific maker
type wasmElementMaker struct {
	client.BaseElementMaker
}

func (c *wasmElementMaker) Get(v any) any {
	return v.(wasm.YWasm)
}
