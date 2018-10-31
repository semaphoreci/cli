package cmd

import (
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a resource.",
	Long:  ``,
}

func init() {
	debugCmd.AddCommand(NewDebugProjectCmd())
	debugCmd.AddCommand(DebugJobCmd)

	RootCmd.AddCommand(debugCmd)
}
