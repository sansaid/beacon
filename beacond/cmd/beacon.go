package cmd

import (
	"beacon/beacond/oci"
	"beacon/beacond/registry"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
)

var Beacon beaconManager

const (
	Probing  ProbeStatus = "probing"
	Outdated ProbeStatus = "outdated"
	Starting ProbeStatus = "starting"
	Exited   ProbeStatus = "exited"
)

type ProbeStatus string

type beacon struct {
	OCIClient      oci.OCIRuntime
	RegistryClient registry.Registry
	close          chan struct{}
	confirmClosing chan struct{}
	Probes         map[string]*Probe
	CleanOnExit    bool
}

type Probe struct {
	close          chan struct{}
	confirmClosing chan struct{}
	resume         chan struct{}
	Namespace      string
	Repo           string
	Status         ProbeStatus
	RunCommand     []string
	CurrentDigest  string
	LatestDigest   string
	LastChecked    time.Time
	LastUpdated    time.Time
}

type beaconManager interface {
	Close()
	ConfirmClosing()
	Start()
	Registry() registry.Registry
	StartProbe(string, string, time.Duration) error
	StopProbe(string, time.Duration) error
	StopProbes(time.Duration) error
	StopManagedContainers(time.Duration) error
}

func NewProbe(namespace string, repo string) *Probe {
	return &Probe{
		Namespace:      namespace,
		Repo:           repo,
		Status:         Starting,
		close:          make(chan struct{}),
		confirmClosing: make(chan struct{}),
		resume:         make(chan struct{}),
	}
}

func NewBeacon(ociClient oci.OCIRuntime, registryClient registry.Registry, cleanOnExit bool) beaconManager {
	if Beacon == nil {
		Beacon = &beacon{
			OCIClient:      ociClient,
			RegistryClient: registryClient,
			CleanOnExit:    cleanOnExit,
			Probes:         make(map[string]*Probe),
			close:          make(chan struct{}),
			confirmClosing: make(chan struct{}),
		}
	}

	return Beacon
}

func (b *beacon) Registry() registry.Registry {
	return b.RegistryClient
}

func (b *beacon) Close() {
	b.close <- struct{}{}
}

func (b *beacon) ConfirmClosing() {
	b.confirmClosing <- struct{}{}
}

func (b *beacon) Start() {
	defer b.ConfirmClosing()

	select {
	case <-b.close:
		if b.CleanOnExit {
			err := b.StopManagedContainers(30 * time.Second)

			if err != nil {
				log.Error(err)
			}
		}

		return
	default:
		for _, probe := range b.Probes {
			if probe.Status == Outdated {
				imageRef := fmt.Sprintf("%s/%s@%s", probe.Namespace, probe.Repo, probe.LatestDigest)

				// Check that a container for this image isn't already running - this can happen if Beacon unexpectedly exits
				// and leaves containers it manages running
				runningContainers, err := b.OCIClient.ContainersUsingImage(imageRef, []string{"running"})

				if err != nil {
					log.Errorf("error checking if image %s is already running:", imageRef, err)
					continue
				}

				if len(runningContainers) > 0 {
					probe.CurrentDigest = probe.LatestDigest
					probe.Resume()
					continue
				}

				err = b.OCIClient.PullImage(imageRef)

				if err != nil {
					log.Errorf("error pulling image %s:", imageRef, err)
					continue
				}

				b.OCIClient.StopContainersByImage(imageRef)
				err = b.OCIClient.RunImage(imageRef)

				if err != nil {
					log.Errorf("error running image %s:", imageRef, err)
					continue
				}

				probe.CurrentDigest = probe.LatestDigest
				probe.Resume()
			}
		}
	}
}

func (b *beacon) StartProbe(namespace string, repo string, delay time.Duration) error {
	probeRef := fmt.Sprintf("%s/%s", namespace, repo)
	b.Probes[probeRef] = NewProbe(namespace, repo)

	go runProbe(b.Probes[probeRef], b.RegistryClient, delay)

	return nil
}

func (b *beacon) StopProbes(delay time.Duration) error {
	select {
	case <-time.Tick(delay):
		return fmt.Errorf("timed out stopping probes")
	default:
		for _, probe := range b.Probes {
			probe.Close()
		}

		b.Probes = make(map[string]*Probe)

		return nil
	}
}

func (b *beacon) StopManagedContainers(delay time.Duration) error {
	select {
	case <-time.Tick(delay):
		return fmt.Errorf("timed out stopping managed containers")
	default:
		for _, probe := range b.Probes {
			imageRef := fmt.Sprintf("%s/%s@%s", probe.Namespace, probe.Repo, probe.LatestDigest)
			err := b.OCIClient.StopContainersByImage(imageRef)

			if err != nil {
				log.Error(err.Error())
			}
		}

		return nil
	}
}

func (b *beacon) StopProbe(probeRef string, delay time.Duration) error {
	b.Probes[probeRef].Close()
	<-b.Probes[probeRef].confirmClosing
	delete(b.Probes, probeRef)

	return nil
}

func (p *Probe) Close() {
	p.close <- struct{}{}
}

func (p *Probe) ConfirmClosing() {
	p.confirmClosing <- struct{}{}
}

func (p *Probe) Resume() {
	p.resume <- struct{}{}
}

func runProbe(probe *Probe, registryClient registry.Registry, delay time.Duration) {
	defer probe.ConfirmClosing()

	prober := func() {
		if probe.Status == Probing {
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

	probe.Status = Probing

	select {
	case <-probe.close:
		return
	case <-probe.resume:
		probe.Status = Probing
		prober()
	default:
		prober()
	}
}
