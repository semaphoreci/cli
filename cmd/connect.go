package cmd

import (
	"fmt"
	"strings"

	"github.com/renderedtext/sem/config"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a Semaphore endpoint",
	Args:  cobra.ExactArgs(2),
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		token := args[1]

		name := strings.Replace(host, ".", "_", -1)

		config.SetActiveContext(name)
		config.SetAuth(token)
		config.SetHost(host)

		fmt.Printf("connected to %s\n", host)
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
