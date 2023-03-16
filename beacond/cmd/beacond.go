package cmd

import (
	"fmt"

	"beacon/beacond/oci"
	"beacon/beacond/registry"

	"github.com/spf13/cobra"
)

var flagBeacondMode enumerable = enumerable{
	allowedValues: []string{"solo", "fleet"},
	currValue:     "solo",
}

var flagOCIRuntime enumerable = enumerable{
	allowedValues: []string{"podman", "docker"},
	currValue:     "podman",
}

var flagRegistry enumerable = enumerable{
	allowedValues: []string{"docker"},
	currValue:     "docker",
}

var flagBeacondPort int
var flagBeacondCleanOnExit bool

var beacond = &cobra.Command{
	Use:   "beacond",
	Short: "beacond is the daemon component that is responsible for running services on your device",
	Run:   beacondHndlr,
}

type enumerable struct {
	allowedValues []string
	currValue     string
}

func (e *enumerable) String() string {
	return e.currValue
}

func (e *enumerable) Set(inputVal string) error {
	for _, a := range e.allowedValues {
		if inputVal == a {
			e.currValue = inputVal
			return nil
		}
	}

	return fmt.Errorf("must be one of: %v", e.allowedValues)
}

func (e *enumerable) Type() string {
	return fmt.Sprintf("%v", e.allowedValues)
}

func init() {
	beacond.PersistentFlags().VarP(&flagBeacondMode, "mode", "m", "The mode to run beacond in")
	beacond.PersistentFlags().VarP(&flagOCIRuntime, "runtime", "r", "The OCI runtime to use")
	beacond.PersistentFlags().VarP(&flagRegistry, "registry", "c", "The container registry to use")
	beacond.PersistentFlags().IntVarP(&flagBeacondPort, "port", "p", 1323, "The port to listen on for commands")
	beacond.PersistentFlags().BoolVar(&flagBeacondCleanOnExit, "clean-up", false, "When beacond exits, whether to also stop containers managed by it")
}

func beacondHndlr(cmd *cobra.Command, args []string) {
	ociClient, err := oci.NewOCIClient(oci.OCIRuntimeType(flagOCIRuntime.currValue))

	if err != nil {
		panic(err)
	}

	registryClient, err := registry.NewRegistry(registry.RegistryType(flagRegistry.currValue))

	if err != nil {
		panic(err)
	}

	beacon := NewBeacon(ociClient, registryClient, flagBeacondCleanOnExit)

	run(beacon, flagBeacondPort)
}

func Execute() error {
	return beacond.Execute()
}
