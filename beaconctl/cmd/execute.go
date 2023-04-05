package cmd

import (
	"github.com/spf13/cobra"
)

var beaconctl = &cobra.Command{
	Use:   "beaconctl",
	Short: "beaconctl is your CLI based controller for beacond",
	Args:  cobra.MinimumNArgs(1),
}

var probe = &cobra.Command{
	Use:   "probe",
	Short: "manage your probes",
	Args:  cobra.MinimumNArgs(1),
}

var probeAdd = &cobra.Command{
	Use:   "add",
	Short: "add a resource",
}

var flagProbe string

func init() {
	// TODO: create a beacond client for beaconctl to communicate with
	probeAdd.PersistentFlags().VarP(&flagProbe, "")
	probe.AddCommand(probeAdd)
	beaconctl.AddCommand(probe)
}

func beaconctlHndlr(cmd *cobra.Command, args []string) {
	panic("not implemented")
}

func Execute() error {
	return beaconctl.Execute()
}
