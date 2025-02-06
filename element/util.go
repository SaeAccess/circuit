package element

import (
	"encoding/json"
	"errors"
	"os/exec"
)

const StdBufferLen = 32e3

func ResolveExe(exe string, cmd string) (path string, err error) {
	path, err = exec.LookPath(exe)
	if err != nil {
		return "", err
	}

	err = exec.Command(path, cmd).Run()
	return path, err
}

func ParseJSONArray[T any](buf []byte) (_ []T, err error) {
	var a []T
	if err = json.Unmarshal(buf, &a); err != nil {
		return nil, err
	}

	return a, nil
}

func ParseJSONArrayFirst[T any](buf []byte) (_ T, err error) {
	data, err := ParseJSONArray[T](buf)
	if err != nil {
		var e T
		return e, err
	}

	if len(data) != 1 {
		var e T
		return e, errors.New("slice does not have one element")
	}

	return data[0], nil
}
