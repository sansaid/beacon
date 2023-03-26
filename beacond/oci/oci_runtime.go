package oci

import (
	"fmt"
)

const (
	Docker OCIRuntimeType = "docker"
	Podman OCIRuntimeType = "podman"
)

type OCIRuntimeType string

type OCIRuntime interface {
	Type() string
	CheckExists() (bool, error)
	PullImage(string) error
	RemoveImages(string, string) error
	RunImage(string) error
	ContainersUsingImage(string, []string) ([]string, error)
	StopContainersByImage(string) error
}

func NewOCIClient(runtime OCIRuntimeType) (OCIRuntime, error) {
	switch runtime {
	case Podman:
		return NewPodman()
	default:
		return nil, fmt.Errorf("runtime not supported: %s", runtime)
	}
}
