package oci

import (
	"os/exec"
)

type PowershellRunner struct{}

func (p PowershellRunner) run(cmds ...string) ([]byte, error) {
	cmd := exec.Command("powershell", cmds...)

	return cmd.CombinedOutput()
}
