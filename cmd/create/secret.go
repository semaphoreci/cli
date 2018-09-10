package cmd_create

import (
	"fmt"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd"
	"github.com/renderedtext/sem/cmd/utils"

	"github.com/spf13/cobra"
)

var createSecretCmd = &cobra.Command{
	Use:   "create secret [NAME]",
	Short: "Create a secret.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		RunCreateSecret(cmd, args)
	},
}

func init() {
	cmd.CreateCmd.AddCommand(createSecretCmd)
}

func RunCreateSecret(cmd *cobra.Command, args []string) {
	name := args[0]

	secret := client.InitSecret(name)
	err := secret.Create()

	utils.Check(err)

	fmt.Printf("Secret '%s' created.", secret.Metadata.Name)
}
