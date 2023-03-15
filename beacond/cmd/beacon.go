package cmd

import (
	"beacon/beacond/oci"
	"beacon/beacond/registry"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

const (
	Pulling  ProbeStatus = "pulling"
	Probing  ProbeStatus = "probing"
	Outdated ProbeStatus = "outdated"
	Starting ProbeStatus = "starting"
	Exited   ProbeStatus = "exited"
)

type ProbeStatus string

type Beacon struct {
	OCIClient      *oci.OCIRuntime
	RegistryClient *registry.Registry
	Probes         map[string]*Probe
}

type Probe struct {
	SignalDone    chan struct{}
	SignalClose   chan struct{}
	Namespace     string
	Repo          string
	Status        ProbeStatus
	RunCommand    []string
	CurrentDigest string
	LatestDigest  string
	LastChecked   time.Time
	LastUpdated   time.Time
}

func NewProbe(namespace string, repo string) *Probe {
	return &Probe{
		Namespace:   namespace,
		Repo:        repo,
		SignalDone:  make(chan struct{}),
		SignalClose: make(chan struct{}),
		Status:      Starting,
	}
}

func NewBeacon(ociClient oci.OCIRuntime, registryClient registry.Registry) *Beacon {
	return &Beacon{
		OCIClient:      &ociClient,
		RegistryClient: &registryClient,
		Probes:         make(map[string]*Probe),
	}
}

func (b *Beacon) RunProbe(namespace string, repo string, delay time.Duration) error {
	probeRef := fmt.Sprintf("%s/%s", namespace, repo)
	b.Probes[probeRef] = NewProbe(namespace, repo)

	go runProbe(b.Probes[probeRef], *b.RegistryClient, delay)

	return nil
}

func (b *Beacon) StopProbes(delay time.Duration) error {
	select {
	case <-time.Tick(delay):
		return fmt.Errorf("timed out stopping probes")
	default:
		for _, probe := range b.Probes {
			probe.SignalClose <- struct{}{}
		}

		b.Probes = make(map[string]*Probe)

		return nil
	}
}

func (b *Beacon) StopProbe(probeRef string, delay time.Duration) error {
	b.Probes[probeRef].SignalClose <- struct{}{}
	<-b.Probes[probeRef].SignalDone
	delete(b.Probes, probeRef)

	return nil
}

func runProbe(probe *Probe, registryClient registry.Registry, delay time.Duration) {
	defer func() { probe.SignalDone <- struct{}{} }()

	probe.Status = Probing

	select {
	case <-probe.SignalClose:
		return
	default:
		if probe.Status != Outdated {
			digest, err := registryClient.LatestImageDigest(probe.Namespace, probe.Repo)

			if err != nil {
				log.Errorf("failed to get latest digest while probing: %s", err)
				probe.Status = Exited
				return
			}

			probe.LastChecked = time.Now()

			if digest != probe.CurrentDigest {
				probe.LatestDigest = digest
				probe.LastUpdated = time.Now()
				probe.Status = Outdated
			}

			time.Sleep(delay)
		}
	}
}
