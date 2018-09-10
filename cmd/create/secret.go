package cmd_create

import (
	"fmt"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/spf13/cobra"
)

var CreateSecretCmd = &cobra.Command{
	Use:   "secret [NAME]",
	Short: "Create a secret.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),

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
