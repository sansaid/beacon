package oci

import "os/exec"

type PosixRunner struct{}

func (p PosixRunner) run(cmds ...string) ([]byte, error) {
	cmd := exec.Command(cmds[0], cmds[1:]...)

	return cmd.CombinedOutput()
}
