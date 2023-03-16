package oci

import (
	"context"
	"os/exec"

	"github.com/labstack/gommon/log"
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

func (p PodmanLinuxClient) PullImage(ref string) error {
	return pullImage(runShell, ref)
}

func (p PodmanLinuxClient) RemoveImage(refPrefix string, olderThanRef string) error {
	images, err := getImages(runShell, refPrefix, olderThanRef, true)

	if err != nil {
		return err
	}

	for _, image := range images {
		err = removeImage(runShell, image)

		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (p PodmanLinuxClient) RunImage(ref string) error {
	return runImage(runShell, ref)
}

func (p PodmanLinuxClient) ContainersUsingImage(ref string, statuses []string) ([]string, error) {
	return containersUsingImage(runShell, ref, statuses)
}

func (p PodmanLinuxClient) StopContainersByImage(ref string) error {
	return stopContainersByImage(runShell, ref)
}

func (p PodmanLinuxClient) CheckExists() (bool, error) {
	return checkExists(runShell)
}
