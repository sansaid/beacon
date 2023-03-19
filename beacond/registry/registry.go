package registry

import "fmt"

const (
	Docker RegistryType = "docker"
)

type RegistryType string

type Registry interface {
	LatestImageDigest(string, string) (string, error)
	TestRepo(string, string) (int, error)
}

func NewRegistry(registryType RegistryType) (Registry, error) {
	switch registryType {
	case Docker:
		return NewDockerRegistry("https://hub.docker.com")
	default:
		return nil, fmt.Errorf("registry type not supported: %s", registryType)
	}
}
