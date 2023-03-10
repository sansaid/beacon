package cmd

import (
	"fmt"

	"beacon/beacond/oci"

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

var flagBeacondPort int

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
	beacond.PersistentFlags().IntVarP(&flagBeacondPort, "port", "p", 1323, "The port to listen on for commands")
}

func beacondHndlr(cmd *cobra.Command, args []string) {
	client, err := oci.NewOCIClient(oci.OCIRuntime(flagOCIRuntime.currValue))

	if err != nil {
		panic(err)
	}

	run(client, flagBeacondPort)
}

func Execute() error {
	return beacond.Execute()
}
