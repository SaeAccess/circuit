package makers

import (
	"reflect"

	"github.com/gocircuit/circuit/client"
	"github.com/gocircuit/circuit/element/wasm"
)

var WasmType = reflect.TypeOf((*wasm.Wasm)(nil)).Elem()

func init() {
	client.RegisterElementMaker(&wasmElementMaker{
		client.NewBaseElementMaker("wasm", WasmType),
	})
}

// implementation for a specific maker
type wasmElementMaker struct {
	client.BaseElementMaker
}
