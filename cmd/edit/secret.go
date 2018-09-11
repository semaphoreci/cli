package cmd_edit

import (
	"fmt"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/handler"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/spf13/cobra"
)

var EditSecretCmd = &cobra.Command{
	Use:     "secret [name]",
	Short:   "Edit a secret.",
	Long:    ``,
	Aliases: []string{"secrets"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		secret, err := client.GetSecret(name)

		utils.Check(err)

		content, err := secret.ToYaml()

		utils.Check(err)

		new_content, err := handler.EditYamlInEditor(secret.ObjectName(), string(content))

		utils.Check(err)

		updated_secret, err := client.InitSecretFromYaml([]byte(new_content))

		utils.Check(err)

		err = updated_secret.Update()

		utils.Check(err)

		fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
	},
}
