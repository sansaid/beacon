package oci

import (
	"fmt"
	"runtime"
)

func NewPodman() (OCIRuntimeAPI, error) {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsPodman()
	case "linux":
		return NewLinuxPodman()
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
