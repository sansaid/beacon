package oci

import (
	"fmt"
	"os"
	"runtime"
)

type OSConfig struct {
	socketRoot string
	socketDir  string
}

func InitOSConfig() (OSConfig, error) {
	switch runtime.GOOS {
	case "linux":
		return OSConfig{
			socketRoot: "unix:",
			socketDir:  os.Getenv("XDG_RUNTIME_DIR"),
		}, nil
	default:
		return OSConfig{}, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
