package cmd

import (
	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client/makers"
	"github.com/gocircuit/circuit/client/podman"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		{
			Name:    "volume",
			Aliases: []string{"v"},
			Usage:   "volume commands",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "create a new volume",
					Args:      true,
					ArgsUsage: "anchor [name]",
					Action:    createVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the network anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "exists",
					Usage:     "Checks if volume exists in local storage",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    checkVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "export",
					Usage:     "exports contents of volume to external tar archive",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    exportVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "import",
					Usage:     "imports contents into a podman volume from specified source",
					Args:      true,
					ArgsUsage: "anchor [source]",
					Action:    importVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "inspect",
					Usage:     "display information for the specified volume hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    inspectVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "mount",
					Usage:     "mount the volume hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    mountVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "reload",
					Usage:     "reload volume",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    reloadVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "scrub",
					Usage:     "abort and remove the container element hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    scrb,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "unmount",
					Usage:     "umount the volume hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    umountVolume,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},

				// other sub commands for podman
			},
		},
	}
	RegisterCommand(cmds...)
}

func createVolume(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	c := dial(x)
	var w []string
	if w, err = ParseAnchor(x); err != nil {
		return err
	}

	var opts *podman.VolumeCreateOptions
	if opts, err = readStdin[podman.VolumeCreateOptions](); err != nil {
		return err
	}

	if x.Args().Len() == 2 {
		opts.Name = x.Args().Get(1)
	}

	if x.Bool("scrub") {
		opts.Scrub = true
	}

	_, err = c.Walk(w).Make(makers.NetworkType, opts)
	if err != nil {
		return errors.Wrapf(err, "createVolume error: %s", err)
	}

	return
}

// TODO should return bool
func checkVolume(x *cli.Context) error {
	return checkAnchorType[podman.Volume](x, anchor.Volume)
}

func exportVolume(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var vol podman.Volume
	if vol, err = getAnchorType[podman.Volume](x, anchor.Volume); err != nil {
		return
	}

	var opts *podman.VolumeExportOptions
	if opts, err = readStdin[podman.VolumeExportOptions](); err != nil {
		return err
	}

	if err = vol.Export(*opts); err != nil {
		return errors.Wrapf(err, "network does not exist")
	}

	return
}

func importVolume(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var vol podman.Volume
	if vol, err = getAnchorType[podman.Volume](x, anchor.Volume); err != nil {
		return
	}

	// source
	var src string
	if x.Args().Len() == 2 {
		src = x.Args().Get(1)

	}

	if err = vol.Import(src); err != nil {
		return errors.Wrapf(err, "network does not exist")
	}

	return
}

func inspectVolume(x *cli.Context) error {
	return inspectAnchorType[podman.Volume](x, anchor.Volume)
}

func mountVolume(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var vol podman.Volume
	if vol, err = getAnchorType[podman.Volume](x, anchor.Volume); err != nil {
		return
	}

	if err = vol.Mount(); err != nil {
		return errors.Wrapf(err, "mount podman volume error")
	}

	return

}

func reloadVolume(x *cli.Context) error {
	return reloadAnchorType[podman.Volume](x, anchor.Volume)
}

func umountVolume(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var vol podman.Volume
	if vol, err = getAnchorType[podman.Volume](x, anchor.Volume); err != nil {
		return
	}

	if err = vol.Unmount(); err != nil {
		return errors.Wrapf(err, "unmount podman volume error")
	}

	return
}
