package podman

type Volume interface {
	// Exists
	Exists() error

	Export(opts VolumeExportOptions) error

	Import(source string) error

	// Inspect the volume
	Inspect() (*InspectVolumeData, error)

	// Mount the volume
	Mount() error

	PeekBytes() []byte

	// Reload the volume from volume plugins
	Reload() error

	// Remove the volume
	Scrub()

	// Unmount the volume
	Unmount() error
}

type VolumeCreateOptions struct {
	Driver  string   `json:"driver,omitempty"`
	Ignore  bool     `json:"ignore,omitempty"`
	Label   []string `json:"label,omitempty"`
	Options []string `json:"options,omitempty"`
	Name    string   `json:"name,omitempty"`

	Scrub bool `json:"scrub,omitempty"`
}

func (c *VolumeCreateOptions) CmdLine() []string {
	var args = []string{"volume", "create"}

	args = appendS(args, "--driver", c.Driver)
	args = appendB(args, "--ignore", c.Ignore)
	args = appendSA(args, "--label", c.Label)
	args = appendSA(args, "--opt", c.Options)

	args = append(args, c.Name)

	return args
}

type VolumeExportOptions struct {
	Output string `json:"output,omitempty"`
}

func (c *VolumeExportOptions) CmdLine(name string) []string {
	var args = []string{"volume", "export"}

	args = appendS(args, "--output", c.Output)
	args = append(args, name)

	return args
}

type VolumeRemoveOptions struct {
	Force bool `json:"force,omitempty"`
	Time  int  `json:"time,omitempty"`
}

func (c *VolumeRemoveOptions) CmdLine(name string) []string {
	var args = []string{"volume", "remove"}

	args = appendB(args, "--force", c.Force)
	args = appendI(args, "--time", c.Time)
	args = append(args, name)

	return args
}
