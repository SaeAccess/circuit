package cmd

import "github.com/urfave/cli/v2"

type Opts func(*cli.App)

var (
	app *cli.App = &cli.App{
		Name:  "circuit",
		Usage: "Circuit server and client tool",
	}
)

// RegisterCommand allows cli application commands to be registered before
// running the application
func RegisterCommand(cmd ...*cli.Command) {
	app.Commands = append(app.Commands, cmd...)
}

// Run runs the circuit cli application
func Run(args []string, opts ...Opts) error {
	for _, o := range opts {
		o(app)
	}

	return app.Run(args)
}

func App() *cli.App {
	return app
}

func WithName(name string) Opts {
	return func(a *cli.App) {
		a.Name = name
	}
}

func WithUsage(u string) Opts {
	return func(a *cli.App) {
		a.Usage = u
	}
}

func WithVersion(v string) Opts {
	return func(a *cli.App) {
		a.Version = v
	}
}

func WithDescription(d string) Opts {
	return func(a *cli.App) {
		a.Description = d
	}
}

func WithBashCompletion(b bool) Opts {
	return func(a *cli.App) {
		a.EnableBashCompletion = b
	}
}

func WithFlags(flags ...cli.Flag) Opts {
	return func(a *cli.App) {
		a.Flags = append(a.Flags, flags...)
	}
}

// add other opts, before, after functions etc.
