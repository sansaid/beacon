package oci

import (
	"context"
	"os/exec"
)

type PodmanLinuxClient struct {
	ctx context.Context
}

func NewLinuxPodman() (OCIRuntime, error) {
	return PodmanLinuxClient{ctx: context.TODO()}, nil
}

func runShell(cmds ...string) ([]byte, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)

	return cmd.CombinedOutput()
}

func (p PodmanLinuxClient) PullImage(digest string) error {
	panic("not implemented")
}

func (p PodmanLinuxClient) RemoveImage(id string) error {
	// See https://docs.docker.com/engine/reference/commandline/images/#filter
	// E.G.: podman images --filter=reference='localhost/vsc-sansaid.github.io-*' --filter 'before=localhost/vsc-sansaid.github.io-0bd29f8740a1596f66a2caa9011b13d8-uid' --format json
	panic("not implemented")
}

func (p PodmanLinuxClient) RunContainer(id string) error {
	panic("not implemented")
}

func (p PodmanLinuxClient) StoppedContainersUsingImage(id string) ([]string, error) {
	panic("not implemented")
}

func (p PodmanLinuxClient) StopContainer(id string) error {
	panic("not implemented")
}

func (p PodmanLinuxClient) CheckExists() (bool, error) {
	return checkExists(runShell)
}
