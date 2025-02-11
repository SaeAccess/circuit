package main

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
			Name:    "network",
			Aliases: []string{"n"},
			Usage:   "network commands",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "create a new network",
					Args:      true,
					ArgsUsage: "anchor [name]",
					Action:    createNetwork,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the network anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "connect",
					Usage:     "add container to the network",
					Args:      true,
					ArgsUsage: "anchor container",
					Action:    connectNetwork,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "disconnect",
					Usage:     "remove container from the network",
					Args:      true,
					ArgsUsage: "anchor container",
					Action:    disconnectNetwork,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "exists",
					Usage:     "Checks if pod exists in local storage",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    checkNetwork,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "inspect",
					Usage:     "display information for the specified pod hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    inspectNetwork,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "reload",
					Usage:     "reload container networks",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    reloadNetwork,
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
					Name:      "update",
					Usage:     "unpause a previously paused container hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    updateNetwork,
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

func createNetwork(x *cli.Context) (err error) {
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

	var opts *podman.NetworkCreateOptions
	if opts, err = readStdin[podman.NetworkCreateOptions](); err != nil {
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
		return errors.Wrapf(err, "createNetwork error: %s", err)
	}

	return
}

// connectNetork
func connectNetwork(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var net podman.Network
	if net, err = getAnchorType[podman.Network](x, anchor.Network); err != nil {
		return
	}

	// connect options
	var opts *podman.NetworkConnectOptions
	if opts, err = readStdin[podman.NetworkConnectOptions](); err != nil {
		return err
	}

	if err = net.Connect(opts); err != nil {
		return errors.Wrapf(err, "podman network connect error")
	}

	return nil
}

func disconnectNetwork(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var net podman.Network
	if net, err = getAnchorType[podman.Network](x, anchor.Network); err != nil {
		return
	}

	// connect options
	var opts *podman.NetworkDisconnectOptions
	if opts, err = readStdin[podman.NetworkDisconnectOptions](); err != nil {
		return err
	}

	if err = net.Disconnect(opts); err != nil {
		return errors.Wrapf(err, "podman network disconnect error")
	}

	return nil
}

// TODO should return bool
func checkNetwork(x *cli.Context) (err error) {
	return checkAnchorType[podman.Network](x, anchor.Network)
}

func inspectNetwork(x *cli.Context) (err error) {
	return inspectAnchorType[podman.Network](x, anchor.Network)
}

func reloadNetwork(x *cli.Context) (err error) {
	return reloadAnchorType[podman.Network](x, anchor.Network)
}

func updateNetwork(x *cli.Context) (err error) {
	return updateAnchorType[podman.Network](x, anchor.Network)
}
