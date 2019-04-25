package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var VERSION = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
