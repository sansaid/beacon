package oci

import (
	"context"
	"os/exec"
)

type PodmanLinuxClient struct {
	ctx context.Context
}

func NewLinuxPodman() (OCIRuntimeAPI, error) {
	return PodmanLinuxClient{ctx: context.TODO()}, nil
}

func runShell(cmds ...string) ([]byte, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)

	return cmd.CombinedOutput()
}

func (p PodmanLinuxClient) Images() []string {
	panic("not implemented")
}

func (p PodmanLinuxClient) CheckExists() (bool, error) {
	return checkExists(runShell)
}
