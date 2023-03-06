package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flagBeacondMode beacondMode = "solo"
var allowedBeacondModeValues []string = []string{"solo", "fleet"}

var beacond = &cobra.Command{
	Use:   "beacond",
	Short: "beacond is the daemon component that is responsible for running services on your device",
	Run:   beacondHndlr,
}

type beacondMode string

func (b *beacondMode) String() string {
	return string(*b)
}

func (b *beacondMode) Set(inputVal string) error {
	for _, a := range allowedBeacondModeValues {
		if inputVal == a {
			*b = beacondMode(inputVal)
			return nil
		}
	}

	return fmt.Errorf("must be one of: %v", allowedBeacondModeValues)
}

func (b *beacondMode) Type() string {
	return fmt.Sprintf("%v", allowedBeacondModeValues)
}

func init() {
	beacond.PersistentFlags().VarP(&flagBeacondMode, "mode", "m", "The mode to run beacond in")
}

func beacondHndlr(cmd *cobra.Command, args []string) {
	// todo: checkOCIRuntimExists(flagOCIRuntime)
	// todo: startWebServer()
}

func Execute() error {
	return beacond.Execute()
}
