package podman

import (
	"log"

	"github.com/gocircuit/circuit/element"
)

var Path string

func ResolvePodman() (string, error) {
	exe, err := element.ResolveExe("podman", "version")
	if err != nil {
		log.Printf("podman not found, error=%v", err)
		return "", err
	}

	Path = exe
	return Path, nil
}
