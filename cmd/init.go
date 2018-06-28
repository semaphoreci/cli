package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    RunInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
  fmt.Println("init called")
}
