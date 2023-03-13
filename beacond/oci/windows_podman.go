package oci

import (
	"context"
	"os/exec"
)

type PodmanWindowsClient struct {
	ctx context.Context
}

func NewWindowsPodman() (OCIRuntimeAPI, error) {
	return PodmanWindowsClient{ctx: context.TODO()}, nil
}

func runPowershell(cmds ...string) ([]byte, error) {
	cmd := exec.Command("powershell", cmds...)

	return cmd.CombinedOutput()
}

func (p PodmanWindowsClient) Images() []string {
	panic("not implemented")
}

func (p PodmanWindowsClient) CheckExists() (bool, error) {
	return checkExists(runPowershell)
}
