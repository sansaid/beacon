package oci

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/labstack/gommon/log"
)

func NewPodman() (OCIRuntime, error) {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsPodman()
	case "linux":
		return NewLinuxPodman()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
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

func runImage(runner cmdRunner, imageRef string) error {
	output, err := runner("podman", "run", imageRef)

	if err != nil {
		return fmt.Errorf("error checking running podman image. Output was: %s; Error was: %s", output, err)
	}

	return nil
}

func pullImage(runner cmdRunner, imageRef string) error {
	output, err := runner("podman", "pull", imageRef)

	if err != nil {
		return fmt.Errorf("error pulling podman image %s. Output was: %s; Error was: %s", imageRef, output, err)
	}

	return nil
}

func removeImage(runner cmdRunner, imageRef string) error {
	output, err := runner("podman", "rm", "image", imageRef)

	if err != nil {
		return fmt.Errorf("error removing podman image %s. Output was: %s; Error was: %s", imageRef, output, err)
	}

	return nil
}

func stopContainersByImage(runner cmdRunner, imageRef string) error {
	containers, err := containersUsingImage(runner, imageRef, []string{"running"})

	if err != nil {
		return err
	}

	for _, container := range containers {
		err = stopContainer(runShell, container)

		if err != nil {
			log.Error(err.Error())
			continue
		}
	}

	return nil
}

func stopContainer(runner cmdRunner, containerID string) error {
	output, err := runner("podman", "stop", containerID)

	if err != nil {
		return fmt.Errorf("error stopping container %s. Output was: %s; Error was: %s", containerID, output, err)
	}

	return nil
}

// See applicable containers: https://docs.docker.com/engine/reference/commandline/ps/#filter
func containersUsingImage(runner cmdRunner, imageRef string, statuses []string) ([]string, error) {
	//  podman ps --filter=ancestor='docker.io/library/httpd@sha256:e4498843f8684e957e3068546ed930b30d43180e2e8c2579d39d637bd2fe79de' --format json
	args := []string{"podman", "ps", "--format", "json"}
	args = append(args, fmt.Sprintf("--filter=ancestor='%s'", imageRef))

	for _, status := range statuses {
		args = append(args, fmt.Sprintf("--filter=status=%s", status))
	}

	output, err := runner(args...)

	if err != nil {
		return []string{}, fmt.Errorf("error getting containers associated with image %s. Output was: %s; Error was: %s", imageRef, output, err)
	}

	var containers []struct {
		ID string `json:"Id"`
	}

	err = json.Unmarshal(output, &containers)

	if err != nil {
		return []string{}, fmt.Errorf("error parsing containers output for image %s. Output was: %s; Error was: %s", imageRef, output, err)
	}

	var containerIDs []string

	for _, container := range containers {
		containerIDs = append(containerIDs, container.ID)
	}

	return containerIDs, nil
}

func getImages(runner cmdRunner, refPrefix string, olderThanImageRef string, dangling bool) ([]string, error) {
	// See https://docs.docker.com/engine/reference/commandline/images/#filter
	// E.G.: podman images --filter=reference='docker.io/library/httpd' --filter 'before=docker.io/library/httpd@sha256:e4498843f8684e957e3068546ed930b30d43180e2e8c2579d39d637bd2fe79de' --format json
	args := []string{"podman", "images", "--format", "json",
		fmt.Sprintf("--filter=reference='%s'", refPrefix),
		fmt.Sprintf("--filter=before='%s'", olderThanImageRef),
		fmt.Sprintf("--filter=dangling=%t", dangling),
	}

	output, err := runner(args...)

	if err != nil {
		return []string{}, fmt.Errorf("error getting images associated with prefix %s. Output was: %s; Error was: %s", refPrefix, output, err)
	}

	var images []struct {
		ID string `json:"Id"`
	}

	err = json.Unmarshal(output, &images)

	if err != nil {
		return []string{}, fmt.Errorf("error parsing images output for ref prefix %s. Output was: %s; Error was: %s", refPrefix, output, err)
	}

	var imageIDs []string

	for _, image := range images {
		imageIDs = append(imageIDs, image.ID)
	}

	return imageIDs, nil
}
