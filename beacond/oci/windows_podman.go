package oci

import (
	"context"
	"os/exec"

	"github.com/labstack/gommon/log"
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

func (p PodmanWindowsClient) Type() string {
	return string(Podman)
}

func (p PodmanWindowsClient) PullImage(ref string) error {
	return pullImage(runPowershell, ref)
}

func (p PodmanWindowsClient) RemoveImage(refPrefix string, olderThanRef string) error {
	images, err := getImages(runPowershell, refPrefix, olderThanRef, true)

	if err != nil {
		return err
	}

	for _, image := range images {
		err = removeImage(runPowershell, image)

		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (p PodmanWindowsClient) RunImage(ref string) error {
	return runImage(runPowershell, ref)
}

func (p PodmanWindowsClient) ContainersUsingImage(ref string, statuses []string) ([]string, error) {
	return containersUsingImage(runPowershell, ref, statuses)
}

func (p PodmanWindowsClient) StopContainersByImage(ref string) error {
	return stopContainersByImage(runPowershell, ref)
}

func (p PodmanWindowsClient) CheckExists() (bool, error) {
	return checkExists(runPowershell)
}
