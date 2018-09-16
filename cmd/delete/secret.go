package cmd_delete

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var DeleteSecretCmd = &cobra.Command{
	Use:     "secret [NAME]",
	Short:   "Delete a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewSecretV1BetaApi()

		err := c.DeleteSecret(name)

		utils.Check(err)

		fmt.Printf("Secret '%s' deleted.\n", name)
	},
}
