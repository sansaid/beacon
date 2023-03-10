package oci

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type OSConfig struct {
	socketRoot string
	socketDir  string
	socketPath string
}

func InitOSConfig() (OSConfig, error) {
	var osc OSConfig
	var xdgDir string

	switch runtime.GOOS {
	case "linux":
		if xdgDir := os.Getenv("XDG_RUNTIME_DIR"); xdgDir == "" {
			return OSConfig{}, fmt.Errorf("no XDG_RUNTIME_DIR exists - cannot determine podman socket location")
		}

		osc = OSConfig{
			socketRoot: "unix:",
			socketDir:  xdgDir,
		}

		if strings.HasPrefix(xdgDir, "/mnt") {
			osc.socketRoot = ""
		}

		osc.socketPath = fmt.Sprintf("%s%s/podman/podman.sock", osc.socketRoot, osc.socketDir)

		return osc, nil
	default:
		return OSConfig{}, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
