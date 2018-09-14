package cmd

import (
	"fmt"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/config"
	"github.com/spf13/cobra"
)

// contextCmd represents the context command
var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Manage contexts for connecting to Semaphore",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0:
			listContexts()
		case 1:
			setContext(args[0])
		default:
			panic("")
		}
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

func listContexts() {
	contexts, err := config.ContextList()

	utils.Check(err)

	active := config.GetActiveContext()

	for _, ctx := range contexts {
		if ctx == active {
			fmt.Print("* ")
			fmt.Println(ctx)
		} else {
			fmt.Print("  ")
			fmt.Println(ctx)
		}
	}
}

func setContext(name string) {
	contexts, err := config.ContextList()

	utils.Check(err)

	for _, ctx := range contexts {
		if ctx == name {
			config.SetActiveContext(name)

			fmt.Printf("switched to context \"%s\"\n", name)
			return
		}
	}

	utils.Fail(fmt.Sprintf("context \"%s\" does not exists", name))
}
