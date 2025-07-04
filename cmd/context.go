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
	ValidArgsFunction: contextValidArgs,
}

func contextValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	var comps []cobra.Completion

	contexts, err := config.ContextList()

	if err != nil {
		return comps, cobra.ShellCompDirectiveNoFileComp
	}

	if len(contexts) == 0 {
		comps = cobra.AppendActiveHelp(comps, "There don't appear to be any contexts defined. Use `sem connect` to configure a context.")
		return comps, cobra.ShellCompDirectiveNoFileComp
	}

	active := config.GetActiveContext()

	for _, context := range contexts {
		compName := context.Name
		compDesc := context.Host

		if compName == active {
			compDesc = fmt.Sprintf("%s (ACTIVE)", compDesc)
		}

		comp := cobra.CompletionWithDesc(compName, compDesc)
		comps = append(comps, comp)
	}

	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	RootCmd.AddCommand(contextCmd)
}

func listContexts() {
	contexts, err := config.ContextNameList()

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
	contexts, err := config.ContextNameList()

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
