package cmd_edit

import (
	"fmt"

	"github.com/spf13/cobra"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	utils "github.com/semaphoreci/cli/cmd/utils"
)

var EditSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Short:   "Edit a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewSecretV1BetaApi()

		secret, err := c.GetSecret(name)

		utils.Check(err)

		content, err := secret.ToYaml()

		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(secret.ObjectName(), string(content))

		utils.Check(err)

		updated_secret, err := models.NewSecretV1BetaFromYaml([]byte(new_content))

		utils.Check(err)

		secret, err = c.UpdateSecret(updated_secret)

		utils.Check(err)

		fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
	},
}
