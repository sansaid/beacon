package oci

import (
	"fmt"
	"strings"
)

const (
	Docker OCIRuntime = "docker"
	Podman OCIRuntime = "podman"
)

type cmdRunner func(cmds ...string) ([]byte, error)

type OCIRuntime string

type OCIRuntimeAPI interface {
	Images() []string
	CheckExists() (bool, error)
}

func NewOCIClient(runtime OCIRuntime) (OCIRuntimeAPI, error) {
	switch runtime {
	case Podman:
		return NewPodman()
	default:
		return nil, fmt.Errorf("runtime not implemented: %s", runtime)
	}
}

func checkExists(runner cmdRunner) (bool, error) {
	output, err := runner("podman", "--version")

	if err != nil {
		return false, fmt.Errorf("error checking podman exists. Output was: %s; Error was: %s", output, err)
	}

	if strings.Contains(string(output), "version") {
		return true, nil
	}

	return false, fmt.Errorf("could not check if podman is running. It either errored unexpectedly or the output was not recognised. Output was: %s; Error was: %s", string(output), err.Error())
}
