package oci

import (
	"context"
	"fmt"
)

const (
	Docker OCIRuntime = "docker"
	Podman OCIRuntime = "podman"
)

type OCIRuntime string

type OCIRuntimeAPI interface {
	Images() []string
}

func NewOCIClient(runtime OCIRuntime) (OCIRuntimeAPI, error) {
	switch runtime {
	case Podman:
		return NewPodman(context.Background())
	default:
		return nil, fmt.Errorf("runtime not implemented: %s", runtime)
	}
}
