package oci

import (
	"context"

	"github.com/containers/podman/v4/pkg/bindings"
)

type PodmanClient struct {
	ctx context.Context
}

func NewPodman(ctx context.Context) (OCIRuntimeAPI, error) {
	config, err := InitOSConfig()

	if err != nil {
		return PodmanClient{}, err
	}

	sockConn, err := bindings.NewConnection(ctx, config.socketPath)

	if err != nil {
		return PodmanClient{}, err
	}

	return PodmanClient{ctx: sockConn}, nil
}

func (p PodmanClient) Images() []string {
	panic("Not implemented")
}
