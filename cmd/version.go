package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.4.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
