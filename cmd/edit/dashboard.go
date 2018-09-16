package cmd_edit

import (
	"fmt"

	"github.com/spf13/cobra"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	utils "github.com/semaphoreci/cli/cmd/utils"
)

var EditDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Edit a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboards", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewDashboardV1AlphaApi()

		dashboard, err := c.GetDashboard(name)

		utils.Check(err)

		content, err := dashboard.ToYaml()

		utils.Check(err)

		new_content, err := utils.EditYamlInEditor(dashboard.ObjectName(), string(content))

		utils.Check(err)

		updated_dashboard, err := models.NewDashboardV1AlphaFromYaml([]byte(new_content))

		utils.Check(err)

		dashboard, err = c.UpdateDashboard(updated_dashboard)

		utils.Check(err)

		fmt.Printf("Dashboard '%s' updated.\n", dashboard.Metadata.Name)
	},
}
