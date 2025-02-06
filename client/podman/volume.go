package podman

type Volume interface {
	// Exists
	Exists() error

	Export() error

	Import() error

	// Inspect the volume
	Inspect() (*InspectVolumeData, error)

	// Mount the volume
	Mount() error

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
}

func (c *VolumeCreateOptions) CmdLine(name string) []string {
	var args = []string{"volume", "create"}

	return args
}

type VolumeExportOptions struct {
	Output string `json:"output,omitempty"`
}

func (c *VolumeExportOptions) CmdLine(name string) []string {
	var args = []string{"volume", "export"}

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
