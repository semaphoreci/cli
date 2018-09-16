package cmd_create

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var CreateDashboardCmd = &cobra.Command{
	Use:     "dashboard [NAME]",
	Short:   "Create a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboard", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewDashboardV1AlphaApi()

		dash := models.NewDashboardV1Alpha(name)
		_, err := c.CreateDashboard(&dash)

		utils.Check(err)

		fmt.Printf("Secret '%s' created.\n", dash.Metadata.Name)
	},
}
