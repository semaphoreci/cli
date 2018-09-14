package cmd_create

import (
	"fmt"

	"github.com/semaphoreci/cli/client"
	"github.com/semaphoreci/cli/cmd/utils"

	"github.com/spf13/cobra"
)

var CreateSecretCmd = &cobra.Command{
	Use:     "secret [NAME]",
	Short:   "Create a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		RunCreateSecret(cmd, args)
	},
}

func RunCreateSecret(cmd *cobra.Command, args []string) {
	name := args[0]

	secret := client.InitSecret(name)
	err := secret.Create()

	utils.Check(err)

	fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
}
