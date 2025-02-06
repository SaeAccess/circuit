package wasm

import "time"

type Status struct {
	ID      string
	Created time.Time
	Path    string
	Args    []string
	//Config          Config
	//State           State
	Image string
	//NetworkSettings NetworkSettings
	ResolvConfPath string
	HostnamePath   string
	HostsPath      string
	Name           string
	Driver         string
	ExecDriver     string
	Volumes        map[string]string
	VolumesRW      map[string]bool
	//HostConfig      HostConfig

	//

	Version Version
}

// Version indicates the wasm runtime version
type Version struct {
	Major uint
	Minor uint
	Patch uint
}

type Result struct {
	Code     uint
	Category uint
	Err      string
}

type Statistics struct {
	InstCount     uint
	InstPerSecond float64
	TotalCost     uint
}
