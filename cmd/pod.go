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
			Name:    "pod",
			Aliases: []string{"p"},
			Usage:   "pod commands",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "create a new pod",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    createPod,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "clone",
					Usage:     "clone a pod",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    clonePod,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "exists",
					Usage:     "Checks if pod exists in local storage",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    checkPod,
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
					Action:    inspectPod,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "pause",
					Usage:     "pause all running containers in the specified pod hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    pausePod,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "peek",
					Usage:     "display information for the specified pod hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    peek,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "restart",
					Usage:     "restart all the containers in the specified pod at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    restart,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "signal",
					Usage:     "Send a signal to the podman container running at anchor",
					Args:      true,
					ArgsUsage: "anchor signal",
					Action:    sgnl,
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
					Name:      "stop",
					Usage:     "stops the running container element hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    stop,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "unpause",
					Usage:     "unpause a previously paused container hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    unpause,
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

func createPod(x *cli.Context) (err error) {
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

	var opts *podman.PodCreateOptions
	if opts, err = readStdin[podman.PodCreateOptions](); err != nil {
		return err
	}

	if x.Bool("scrub") {
		opts.Scrub = true
	}

	_, err = c.Walk(w).Make(makers.PodType, opts)
	if err != nil {
		return errors.Wrapf(err, "createPod error: %s", err)
	}

	return
}

func clonePod(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	var p podman.Pod
	if p, err = getAnchorType[podman.Pod](x, anchor.Pod); err != nil {
		return
	}

	var opts *podman.PodCloneOptions
	if opts, err = readStdin[podman.PodCloneOptions](); err != nil {
		return err
	}

	if _, err = p.Clone(opts); err != nil {
		return errors.Wrapf(err, "podman pod clone error")
	}

	return
}

// TODO this should return bool, true - exists, false otherwise
func checkPod(x *cli.Context) (err error) {
	return checkAnchorType[podman.Pod](x, anchor.Pod)
}

// Inspect the configuration of the container
func inspectPod(x *cli.Context) (err error) {
	return inspectAnchorType[podman.Pod](x, anchor.Pod)
}

// Pause the processes in the container
func pausePod(x *cli.Context) (err error) {
	return pauseAnchorType[podman.Pod](x, anchor.Pod)
}

func restart(x *cli.Context) (err error) {
	return restartAnchorType[podman.Pod](x, anchor.Pod)
}
