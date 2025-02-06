package wasm

import (
	"encoding/gob"
	"encoding/json"
)

// init register any global types.
func init() {
	gob.Register(Execute{})
}

type LogLevel string

const (
	LOG_ERROR LogLevel = "error"
	LOG_WARN  LogLevel = "warn"
	LOG_DEBUG LogLevel = "debug"
)

// Execute parameterizes the execution of a wasm application module or function
// Provides configuration for runtime objects (VM, WASI, etc)
type Execute struct {
	LogLevel          LogLevel
	Function          Functional
	WasiConf          WasiConf
	AddProposals      Proposals
	RemoveProposals   Proposals
	MaxMemoryPageSize uint
	ForceInterpreter  bool
	CompileOpts       CompileOpts
	StatsOpts         StatisticsOpts

	Scrub bool
}

type Limit struct {
	Min uint
	Max uint
}

type Functional struct {
	IsFunction bool
	Name       string
	Args       []string
}

type WasiConf struct {
	Enable   bool
	Args     []string
	Env      []string
	PreOpens []string
}

// Proposals indicates the WebAssembly proposals to turn on/off
type Proposals struct {
	ImportExportMutGlobals       bool
	NonTrapFloatToIntConversions bool
	SignEextensionOperators      bool
	MultiValue                   bool
	BulkMemoryOperations         bool
	ReferenceTypes               bool
	SIMD                         bool
	TailCall                     bool
	MultiMemories                bool
	Annotations                  bool
	Memory64                     bool
	ExceptionHandling            bool
	ExtendedConst                bool
	Threads                      bool
	FunctionReferences           bool
}

type CompileOpts struct {
	OptLevel      CompilerOptLevel
	OutputFormat  CompilerOutputFormat
	DumpIR        bool
	GenericBinary bool
}

type CompilerOutputFormat string

const (
	NativeOutputFormat CompilerOutputFormat = "native"
	WasmOutputFormat   CompilerOutputFormat = "wasm"
)

type CompilerOptLevel string

const (
	OptLevel_O0 CompilerOptLevel = "O0"
	OptLevel_O1 CompilerOptLevel = "O1"
	OptLevel_O2 CompilerOptLevel = "O2"
	OptLevel_O3 CompilerOptLevel = "O3"
	OptLevel_Os CompilerOptLevel = "Os"
	OptLevel_Oz CompilerOptLevel = "Oz"
)

type StatisticsOpts struct {
	CountInstructions bool
	TimeMeasurement   bool
	CostMeasurement   bool
}

func ParseType[T any](src string) (*T, error) {
	var x T
	if err := json.Unmarshal([]byte(src), &x); err != nil {
		return nil, err
	}
	return &x, nil
}
