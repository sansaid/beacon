package oci

import (
	"context"
	"os/exec"
)

type PodmanWindowsClient struct {
	ctx context.Context
}

func NewWindowsPodman() (OCIRuntime, error) {
	return PodmanWindowsClient{ctx: context.TODO()}, nil
}

func runPowershell(cmds ...string) ([]byte, error) {
	cmd := exec.Command("powershell", cmds...)

	return cmd.CombinedOutput()
}

func (p PodmanWindowsClient) PullImage(digest string) error {
	panic("not implemented")
}

func (p PodmanWindowsClient) RemoveImage(id string) error {
	// See https://docs.docker.com/engine/reference/commandline/images/#filter
	// E.G.: podman images --filter=reference='localhost/vsc-sansaid.github.io-*' --filter 'before=localhost/vsc-sansaid.github.io-0bd29f8740a1596f66a2caa9011b13d8-uid' --format json
	panic("not implemented")
}

func (p PodmanWindowsClient) RunContainer(id string) error {
	panic("not implemented")
}

func (p PodmanWindowsClient) StoppedContainersUsingImage(id string) ([]string, error) {
	panic("not implemented")
}

func (p PodmanWindowsClient) StopContainer(id string) error {
	panic("not implemented")
}

func (p PodmanWindowsClient) CheckExists() (bool, error) {
	return checkExists(runPowershell)
}
