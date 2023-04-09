package cmd

import (
	"github.com/spf13/cobra"
)

var RESOURCES = []string{"probe"}

var beaconctl = &cobra.Command{
	Use:   "beaconctl",
	Short: "beaconctl is your CLI based controller for beacond",
	Args:  cobra.MinimumNArgs(1),
}

var crudCmds = []*cobra.Command{
	{
		Use:       "create",
		Short:     "create a resource",
		ValidArgs: RESOURCES,
	},
	{
		Use:       "delete",
		Short:     "delete a resource",
		ValidArgs: RESOURCES,
	},
	{
		Use:       "list",
		Short:     "list a resource",
		ValidArgs: RESOURCES,
	},
	{
		Use:       "describe",
		Short:     "describe a resource",
		ValidArgs: append(RESOURCES, "beacon"),
	},
}

func init() {
	// TODO: create a beacond client for beaconctl to communicate with
	initialiseCrudCmds()
}

func beaconctlHndlr(cmd *cobra.Command, args []string) {
	panic("not implemented")
}

func initialiseCrudCmds() {
	for _, cmd := range crudCmds {
		beaconctl.AddCommand(cmd)
	}
}

func Execute() error {
	return beaconctl.Execute()
}
