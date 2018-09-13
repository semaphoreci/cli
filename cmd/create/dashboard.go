package cmd_create

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/renderedtext/sem/api"
	"github.com/renderedtext/sem/api/client/semaphore_dashboards_v1alpha_dashboards_api"
	models "github.com/renderedtext/sem/api/models"
	"github.com/renderedtext/sem/cmd/utils"
)

var CreateDashboardCmd = &cobra.Command{
	Use:     "dashboard [NAME]",
	Short:   "Create a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboard", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		d := models.SemaphoreDashboardsV1alphaDashboard{
			Metadata: &models.SemaphoreDashboardsV1alphaDashboardMetadata{
				Name: name,
			},
			Spec: &models.SemaphoreDashboardsV1alphaDashboardSpec{},
		}

		c := api.DefaultClient()
		params := semaphore_dashboards_v1alpha_dashboards_api.NewCreateDashboardParams().WithBody(&d)
		_, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.CreateDashboard(params)

		if err != nil {
			utils.Check(err)
		}

		fmt.Printf("Dashboard '%s' created.\n", d.Metadata.Name)
	},
}
