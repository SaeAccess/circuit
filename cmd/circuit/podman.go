package main

import (
	"fmt"

	"github.com/gocircuit/circuit/anchor"
	"github.com/gocircuit/circuit/client/makers"
	container "github.com/gocircuit/circuit/client/podman"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds := []*cli.Command{
		{
			Name:    "container",
			Aliases: []string{"c"},
			Usage:   "Container commands",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Usage:     "create a new container",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    createContainer,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "run",
					Usage:     "run a container",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    runContainer,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.BoolFlag{Name: "scrub", Usage: "scrub the process anchor automatically on exit"},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "checkpoint",
					Usage:     "Checkpoint a runnning container",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    checkPoint,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "exec",
					Usage:     "Execute the specified command inside the runnning container",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    exec,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "inspect",
					Usage:     "display information for the specified container hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    inspect,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "pause",
					Usage:     "pause the specified container hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    pause,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "peek",
					Usage:     "display information for the specified container hosted at anchor",
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
					Name:      "ports",
					Usage:     "list port mappings for the specified container hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    ports,
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "dial", Aliases: []string{"d"}, Value: "", Usage: "circuit member to dial into"},
						&cli.StringFlag{Name: "discover", Value: "228.8.8.8:8822", Usage: "Multicast address for peer server discovery", EnvVars: []string{"CIRCUIT_DISCOVER"}},
						&cli.StringFlag{Name: "hmac", Value: "", Usage: "File containing HMAC credentials. Use RC4 encryption.", EnvVars: []string{"CIRCUIT_HMAC"}},
					},
				},
				{
					Name:      "restore",
					Usage:     "restores a container from a checkpoint to the specified anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    restore,
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
					Name:      "start",
					Usage:     "start the container element hosted at anchor",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    start,
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
				{
					Name:      "wait",
					Usage:     "block until one or more container stop",
					Args:      true,
					ArgsUsage: "anchor",
					Action:    wait,
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

func createContainer(x *cli.Context) (err error) {
	_, err = makeContainer(x)
	return err
}

func runContainer(x *cli.Context) error {
	conn, err := makeContainer(x)
	if err != nil {
		return err
	}

	// run the container
	return conn.Start()
}

func wrapError(r any) (err error) {
	err, ok := r.(error)
	if ok {
		err = errors.Wrapf(err, "error, likely due to missing server or misspelled anchor: %v", r)
	} else {
		err = errors.Wrapf(fmt.Errorf("error: %v", r), "error, likely due to missing server or misspelled anchor")
	}

	return
}

func makeContainer(x *cli.Context) (con container.Container, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	c := dial(x)
	args := x.Args()
	if args.Len() != 1 {
		return nil, errors.New("makeContainer needs an anchor argument")
	}
	w, _ := parseGlob(args.First())
	var opts *container.ContainerCreateOptions
	if opts, err = readStdin[container.ContainerCreateOptions](); err != nil {
		return nil, err
	}

	if x.Bool("scrub") {
		opts.Scrub = true
	}

	_, err = c.Walk(w).Make(makers.ContainerType, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "makeContainer error: %s", err)
	}

	return
}

func checkPoint(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	// read container checkpoint options from stdin
	var opts *container.ContainerCheckpointOptions
	if opts, err = readStdin[container.ContainerCheckpointOptions](); err != nil {
		return err
	}

	if err = con.CheckPoint(opts); err != nil {
		return errors.Wrapf(err, "podman container checkpoint error: %v", err)
	}
	return
}

// Exec runs a process in the container
func exec(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	// read container checkpoint options from stdin
	var opts *container.ContainerExecOptions
	if opts, err = readStdin[container.ContainerExecOptions](); err != nil {
		return err
	}

	var result []byte
	if result, err = con.Exec(opts); err != nil {
		return errors.Wrapf(err, "podman container exec error: %v", err)
	}

	fmt.Println(string(result))
	return
}

// Inspect the configuration of the container
func inspect(x *cli.Context) (err error) {
	return inspectAnchorType[container.Container](x, anchor.Container)
}

// Pause the processes in the container
func pause(x *cli.Context) (err error) {
	return pauseAnchorType[container.Container](x, anchor.Container)
}

// Stats()
// Peek at the container's resource ussage statistics
// TODO this should be in generic peek
// func podmanPeek(x *cli.Context) error {
// 	return nil

// }

// Get port mappings
func ports(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	p := con.Ports()
	fmt.Printf("%v", p)
	return
}

// Restore the container from a checkpoint
func restore(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	// read container checkpoint options from stdin
	var opts *container.ContainerRestoreOptions
	if opts, err = readStdin[container.ContainerRestoreOptions](); err != nil {
		return err
	}

	if err = con.Restore(opts); err != nil {
		return errors.Wrapf(err, "podman container restore error: %v", err)
	}

	return
}

// RunLabel runs the command specified in the label
func runLabel(x *cli.Context) error {
	return nil
}

// Scrub removes the container
// generic scrb in
// func scrub(x *cli.Context) error {
// }

func start(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	if err = con.Start(); err != nil {
		return errors.Wrapf(err, "podman container start error")
	}

	return

}

// Stop the container
func stop(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	var opts *container.ContainerStopOpts
	if opts, err = readStdin[container.ContainerStopOpts](); err != nil {
		return err
	}

	if err = con.Stop(opts); err != nil {
		return errors.Wrapf(err, "podman container stop error: %v", err)
	}

	return
}

// Unpause the processes in the container
func unpause(x *cli.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = wrapError(r)
		}
	}()

	// get the container anchor
	var con container.Container
	if con, err = getAnchorType[container.Container](x, anchor.Container); err != nil {
		return
	}

	if err = con.Unpause(); err != nil {
		return errors.Wrapf(err, "podman container unpause error: %v", err)
	}

	return
}
