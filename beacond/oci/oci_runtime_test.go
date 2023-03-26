package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOCIClient(t *testing.T) {
	var err error

	_, err = NewOCIClient(Podman)
	assert.NoError(t, err)

	_, err = NewOCIClient(Docker)
	assert.Errorf(t, err, "runtime not supported: %s", string(Docker))
}
