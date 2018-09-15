package cmd_create

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

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
		name := args[0]

		c := client.NewSecretV1BetaApi()

		secret := models.NewSecretV1Beta(name)
		_, err := c.CreateSecret(&secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' created.\n", secret.Metadata.Name)
	},
}
