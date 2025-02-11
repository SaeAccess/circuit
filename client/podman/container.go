package podman

import "io"

// Container is the interface for managing a container
type Container interface {
	// Checkpoint the container
	CheckPoint(opts *ContainerCheckpointOptions) error

	// Clone an existing container
	// Clone(opts *ContainerCloneOptions) (Container, error)

	// Exec runs a process in the container
	Exec(opts *ContainerExecOptions) ([]byte, error)

	// Inspect the configuration of the container
	Inspect() (*InspectContainerData, error)

	IsDone() bool

	// Pause the processes in the container
	Pause() error

	// Stats()
	// Peek at the container's resource ussage statistics
	Peek() (*InspectContainerData, error)
	PeekBytes() []byte

	// Get port mappings
	Ports() []string

	// Restore the container from a checkpoint
	Restore(opts *ContainerRestoreOptions) error

	// RunLabel runs the command specified in the label
	RunLabel() error

	// Scrub removes the container
	Scrub()

	// Start the container
	Start() error

	// Stop the container
	Stop(opts *ContainerStopOpts) error

	// Unpause the processes in the container
	Unpause() error

	// Kill the container with the specified signal
	Signal(sig string) error

	// Wait for the container to exit
	Wait() (*InspectContainerData, error)

	Stdin() io.WriteCloser
	Stdout() io.ReadCloser
	Stderr() io.ReadCloser
}
