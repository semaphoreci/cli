package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ReleaseVersion = "dev"
	ReleaseCommit  = "none"
	ReleaseDate    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v, commit %v, built at %v\n",
			ReleaseVersion,
			ReleaseCommit,
			ReleaseDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
